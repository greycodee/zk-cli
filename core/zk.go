package core

import (
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type ZK struct {
	conn *zk.Conn
}

func NewZKClient(host string) (zkClient *ZK,err error)  {
	if host == "" {
		return nil,errors.New("host is nil")
	}

	// 创建zk连接地址
	hosts := []string{host}
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		return
	}

	z := ZK{}
	z.conn = conn
	zkClient = &z
	return
}


func (z ZK) Children(path string) (list []string, stat *zk.Stat,err error) {
	list, stat, err = z.conn.Children(path)
	return
}

func (z ZK) HaveChildren(path string) bool {
	_, stat, _ := z.conn.Children(path)
	num := stat.NumChildren
	return num > 0
}

func (z ZK) GetData(path string) (data string, stat *zk.Stat,err error) {
	get, stat, err := z.conn.Get(path)
	if err != nil {
		return
	}
	data = string(get)
	return
}
