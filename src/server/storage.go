package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

const SS_PREFIX = "ss:"

type User struct {
	Name     string
	Password string
	Port     int
	Method   string
	Limit    string
	EndDate  time.Time   //到期删除时间
	State 	 string        //账号状态 ok del stop

}

type Storage struct {
	pool *redis.Pool
}

func NewStorage() *Storage {

	pool := redis.NewPool(func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("unix", "/tmp/redis.sock")
			conn.Do("AUTH","%$s%dd$%d#s^df#$a^fd%sf*^&(d*d&^)gh*^jk*e(*&e*s#%")
			return
		}, 3)
	return &Storage{pool}
}

func (s *Storage) Get(key string) (user User, err error) {
	fullkey := SS_PREFIX + key

	return s.get(fullkey)
}

func (s *Storage) get(fullkey string) (user User, err error) {
	var data []byte
	var conn = s.pool.Get()
	defer conn.Close()
	data, err = redis.Bytes(conn.Do("GET", fullkey))
	if err != nil {
		return
	}
	//fmt.Println(string(data))

	err = json.Unmarshal(data, &user)
	//fmt.Println(user)
	return
}

func (s *Storage) GetList() (userList []User, err error) {
	//var data []byte
	var conn = s.pool.Get()
	defer conn.Close()
	key := "user*"
	data, err := redis.Strings(conn.Do("keys", SS_PREFIX+key))

	if err != nil {
		fmt.Printf("GetList()\n", err)
		return
	}

	userList = make([]User, len(data))
	for k, v := range data {
		_user, _ := s.get(v)
		userList[k] = _user
	}

	return
}
func (s *Storage) Set(key string, user User) (err error) {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	conn := s.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", SS_PREFIX+key, data)
	return
}

//将日志信息写入redis
func (s *Storage) Log(key ,data string) (err error) {

	conn := s.pool.Get()
	defer conn.Close()
	_, err = conn.Do("SET", SS_PREFIX+"log:"+key, data)
	return
}
func (s *Storage) IncrSize(key string, incr int) (score int64, err error) {
	var conn = s.pool.Get()
	defer conn.Close()
	score, err = redis.Int64(conn.Do("INCRBY", SS_PREFIX+key, incr))
	return
}

func (s *Storage) GetSize(key string) (score int64, err error) {
	var conn = s.pool.Get()
	defer conn.Close()
	score, err = redis.Int64(conn.Do("GET", SS_PREFIX+key))
	return
}
