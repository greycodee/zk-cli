package core

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type ZK struct {
	Host string
	conn *zk.Conn
}

func (z *ZK) newClient() (zkClient *ZK,err error)  {
	// 创建zk连接地址
	hosts := []string{z.Host}
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		return
	}
	z.conn = conn
	zkClient = z
	return
}
