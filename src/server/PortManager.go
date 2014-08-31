package main

import (
	"net"
	"sync"
	"time"
	"strconv"
	"fmt"
)

var listenerManager = ListenerManager{portListener: map[string]*PortListener{}}

type PortListener struct {
	user     User
	listener net.Listener
}

type ListenerManager struct {
	sync.Mutex
	portListener map[string]*PortListener
}

func (pm *ListenerManager) add(user User , listener net.Listener) {
	pm.Lock()
	pm.portListener[user.Name] = &PortListener{user, listener}
	pm.Unlock()
}

func (pm *ListenerManager) get(user User) (pl *PortListener, ok bool) {
	pm.Lock()
	pl, ok = pm.portListener[user.Name]
	pm.Unlock()
	return
}

func (pm *ListenerManager) del(user User) {
	pl, ok := pm.get(user)
	if !ok {
		return
	}
	pl.listener.Close()
	pm.Lock()
	delete(pm.portListener, user.Name)
	pm.Unlock()
}

//统一管理入口，注册后台监控链接
func (pm *ListenerManager) updateUser() {
	//得到最新的用户列表
	Userlist, _ := storage.GetList()


	//先关闭已删除的账号
	for _, v := range pm.portListener {
		if !isIn(v.user, Userlist) {
			debug.Println("删除用户",v.user)
			pm.del(v.user)
		}
	}



	for _, user := range Userlist {


		switch user.State{
		case "ok":
			if !pm.ishas(user) {
				if _, ok := pm.get(user); !ok {
					debug.Printf("端口第一次注册", user.Name, user.Port)
					//开启新的链接
					go run(user)
				}else{
					debug.Printf("端口已注册,更新", user.Name, user.Port)
					//先终止原链接，清除端口链接跟密码绑定关系
					pm.del(user)
					//开启新的链接
					go run(user)
				}


			}else {
				if _, ok := pm.get(user); !ok {
					//不存在监听端口的话，新开一个端口
					debug.Printf("新加入用户", user.Name, user.Port)
					go run(user)
				}
			}
		default:
			debug.Printf("删除用户", user.Name, user.Port)
			pm.del(user)
		}

	}
}

//查找是否包含
func isIn(user User, userList []User) (flag bool) {
	flag = false
	for _, v := range userList {
		if v == user {
			flag = true
			break
		}
	}
	return
}

//检查用户是否存在于已注册用户中
func (pm *ListenerManager) ishas(user User) (flag bool) {
	flag = false
	for _, v := range pm.portListener {
		if v.user.Name == user.Name && v.user == user {
			flag = true
			break
		}
	}
	return
}

func CheckAuth(user User, ln net.Listener) {
	debug.Println("用户", user.Name, user.Password)

	//超流量关停
	currentSize, _ := storage.GetSize("flow:" + user.Name)
	limit, _ := strconv.Atoi(user.Limit)
	if currentSize >= int64(limit) {
		debug.Println("数据超限，关停:", user.Limit, currentSize, user.Name, user.Password)
		go storage.Log(user.Name, fmt.Sprintf("数据配额%s,已用%s，关停%s:%s", user.Limit, currentSize, user.Name, user.Password))
		ln.Close()
	}

	//超时关停
	if time.Now().After(user.EndDate) {
		debug.Println("账户到期，关停:", user.EndDate, user.Name, user.Password)
		go storage.Log(user.Name, fmt.Sprintf("账户到期%s,关停%s:%s", user.EndDate, user.Name, user.Password))
		ln.Close()
	}
}
