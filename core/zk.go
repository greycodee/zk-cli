package core

import (
	"github.com/greycodee/zk-cli/log"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type ZK struct {
	conn *zk.Conn
	tuiLog *log.Log
}

func NewZKClient(host string, log *log.Log) (zkClient *ZK)  {
	if host == "" {
		panic("zookeeper 地址不能为空")
	}

	// 创建zk连接地址
	hosts := []string{host}
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		panic("连接 zookeeper 失败")
	}

	z := ZK{
		tuiLog: log,
	}
	z.conn = conn
	zkClient = &z
	return
}


func (z ZK) Children(path string) (list []string, stat *zk.Stat) {
	list, stat, err := z.conn.Children(path)
	if err != nil {
		z.tuiLog.Info(err.Error())
	}
	return
}

func (z ZK) HaveChildren(path string) bool {
	_, stat, err := z.conn.Children(path)
	if err != nil {
		z.tuiLog.Info(err.Error())
	}else{
		num := stat.NumChildren
		return num > 0
	}
	return false
}

func (z ZK) GetData(path string) (data string, stat *zk.Stat) {
	get, stat, err := z.conn.Get(path)
	if err != nil {
		z.tuiLog.Info(err.Error())
	}
	data = string(get)
	return
}
