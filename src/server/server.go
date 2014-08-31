package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"github.com/cyfdecyf/leakybuf"
	"time"

)

var storage *Storage

var debug ss.DebugLog

func getRequest(conn *ss.Conn) (host string, extra []byte, err error) {
	const (
		idType  = 0 // address type index
		idIP0   = 1 // ip addres start index
		idDmLen = 1 // domain address length index
		idDm0   = 2 // domain address start index

		typeIPv4 = 1 // type is ipv4 address
		typeDm   = 3 // type is domain address
		typeIPv6 = 4 // type is ipv6 address

		lenIPv4   = 1 + net.IPv4len + 2 // 1addrType + ipv4 + 2port
		lenIPv6   = 1 + net.IPv6len + 2 // 1addrType + ipv6 + 2port
		lenDmBase = 1 + 1 + 2           // 1addrType + 1addrLen + 2port, plus addrLen
	)

	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 256(max length address) + 2(port)
	buf := make([]byte, 260)
	var n int
	// read till we get possible domain length field
	ss.SetReadTimeout(conn)
	if n, err = io.ReadAtLeast(conn, buf, idDmLen+1); err != nil {
		return
	}

	reqLen := -1
	switch buf[idType] {
	case typeIPv4:
		reqLen = lenIPv4
	case typeIPv6:
		reqLen = lenIPv6
	case typeDm:
		reqLen = int(buf[idDmLen])+lenDmBase
	default:
		err = errors.New(fmt.Sprintf("addr type %d not supported", buf[idType]))
		return
	}

	if n < reqLen { // rare case
		ss.SetReadTimeout(conn)
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return
		}
	} else if n > reqLen {
		// it's possible to read more than just the request head
		extra = buf[reqLen:n]
	}

	// Return string for typeIP is not most efficient, but browsers (Chrome,
	// Safari, Firefox) all seems using typeDm exclusively. So this is not a
	// big problem.
	switch buf[idType] {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0 + buf[idDmLen]])
	}
	// parse port
	port := binary.BigEndian.Uint16(buf[reqLen - 2 : reqLen])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

const logCntDelta = 100

var connCnt int
var nextLogConnCnt int = logCntDelta

///得到数据大小 by weibh 2014 08 26

func LogSize(user User , size int) {
	go storage.IncrSize("flow:"+user.Name, size)
	go debug.Println("本次数据大小:", user.Name, size)
}

var readTimeout time.Duration

const bufSize = 4096
const nBuf = 2048
const (
	NO_TIMEOUT = iota
	SET_TIMEOUT
)

func SetReadTimeout(c net.Conn) {
	if readTimeout != 0 {
		c.SetReadDeadline(time.Now().Add(readTimeout))
	}
}

var pipeBuf = leakybuf.NewLeakyBuf(nBuf, bufSize)

// PipeThenClose copies data from src to dst, closes dst when done.
func PipeThenClose(src, dst net.Conn, timeoutOpt int, user User) {

	defer dst.Close()
	buf := pipeBuf.Get()
	defer pipeBuf.Put(buf)
	for {
		if timeoutOpt == SET_TIMEOUT {
			SetReadTimeout(src)
		}
		n, err := src.Read(buf)
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		//		if n > 0 {
		//			if _, err = dst.Write(buf[0:n]); err != nil {
		//				log.Println("write:", err)
		//				break
		//			}
		//		}

		/////add by weibh 2014 08 26

		if n > 0 {
			size, err := dst.Write(buf[0:n])
			if err != nil {
				log.Println("write:", err)
				break
			}
			//记录大小
			go LogSize(user, size)
			//
			//			if total_size > user.Limit {
			//				return
			//			}
		}
		if err != nil {
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			/*
				if bool(Debug) && err != io.EOF {
					Debug.Println("read:", err)
				}
			*/
			break
		}
	}
}

func handleConnection(conn *ss.Conn, user User) {
	var host string

	connCnt++ // this maybe not accurate, but should be enough
	if connCnt-nextLogConnCnt >= 0 {
		// XXX There's no xadd in the atomic package, so it's difficult to log
		// the message only once with low cost. Also note nextLogConnCnt maybe
		// added twice for current peak connection number level.
		log.Printf("Number of client connections reaches %d\n", nextLogConnCnt)
		nextLogConnCnt += logCntDelta
	}

	// function arguments are always evaluated, so surround debug statement
	// with if statement
	if debug {
		debug.Printf("new client %s->%s\n", conn.RemoteAddr().String(), conn.LocalAddr())
	}
	closed := false
	defer func() {
		//fmt.Printf("host:%s,RemoteAddr:%s,LocalAddr:%s\n",host, conn.RemoteAddr().String(),conn.LocalAddr())

		if debug {
			debug.Printf("closed pipe %s<->%s\n", conn.RemoteAddr(), host)
		}
		connCnt--
		if !closed {
			conn.Close()
		}
	}()

	host, extra, err := getRequest(conn)
	if err != nil {
		debug.Println("error getting request", conn.RemoteAddr(), conn.LocalAddr(), err)
		return
	}
	debug.Println("connecting", host)
	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			log.Println("dial error:", err)
		} else {
			debug.Println("error connecting to:", host, err)
		}
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()
	// write extra bytes read from
	if extra != nil {
		debug.Println("getRequest read extra data, writing to remote, len", len(extra))
		if _, err = remote.Write(extra); err != nil {
			debug.Println("write request extra error:", err)
			return
		}

		///得到数据大小 by weibh 2014 08 26
		//		size, err := remote.Write(extra)
		//		//额外验证信息
		//		go LogSize(user, size)
		//
		//		if err != nil {
		//			debug.Println("write request extra error:", err)
		//			return
		//		}
		///得到数据大小 by weibh 2014 08 26


	}
	if debug {
		debug.Printf("piping %s<->%s", conn.RemoteAddr(), host)
	}
	//	go ss.PipeThenClose(conn, remote, ss.SET_TIMEOUT)
	//	ss.PipeThenClose(remote, conn, ss.NO_TIMEOUT)

	//add by weibh 20140826 增加对user 进行流量统计
	go PipeThenClose(conn, remote, ss.SET_TIMEOUT, user)
	PipeThenClose(remote, conn, ss.NO_TIMEOUT, user)
	closed = true
	return
}

//等待动态更新新号
//ps -ef|grep ./server|grep -v grep|awk '{printf $2}'|xargs kill -1
//kill -1 发送SIGHUP信号
func waitSignal() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	for sig := range sigChan {
		if sig == syscall.SIGHUP {
			fmt.Println("收到信号")
			listenerManager.updateUser()
		} else {
			// is this going to happen?
			log.Printf("caught signal %v, exit", sig)
			os.Exit(0)
		}
	}
}

//func run(port, password string) {
func run(user User) {
	debug.Println("已注册", user)
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(user.Port))
	if err != nil {
		log.Printf("error listening port %v: %v\n", user.Port, err)
		return
	}
	log.Printf("server listening port %v ...\n", user.Port)

	listenerManager.add(user, ln)
	var cipher *ss.Cipher

	for {
		conn, err := ln.Accept()
		if err != nil {
			// listener maybe closed to update password
			debug.Printf("accept error: %v\n", err)
			return
		}

		go CheckAuth(user, ln)


		// Creating cipher upon first connection.
		if cipher == nil {
			debug.Println("creating cipher for port:", user.Port)
			//cipher, err = ss.NewCipher(config.Method, user.Password)
			cipher, err = ss.NewCipher(user.Method, user.Password)
			if err != nil {
				debug.Printf("Error generating cipher for port: %s %v\n", user.Port, err)
				conn.Close()
				continue
			}
		}
		//add password by weibh 2014 08 26
		go handleConnection(ss.NewConn(conn, cipher.Copy()), user)
	}
}

var config *ss.Config

func main() {
	log.SetOutput(os.Stdout)

	var printVer bool
	var core int

	flag.BoolVar(&printVer, "version", false, "print version")

	flag.IntVar(&core, "core", 0, "maximum number of CPU cores to use, default is determinied by Go runtime")
	flag.BoolVar((*bool)(&debug), "d", false, "print debug message")

	flag.Parse()

	if printVer {
		ss.PrintVersion()
		os.Exit(0)
	}

	ss.SetDebug(debug)



	storage = NewStorage()

	if core > 0 {
		runtime.GOMAXPROCS(core)
	}

	//在这里统一管理端口
	listenerManager.updateUser()
	waitSignal()
}
