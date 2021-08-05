package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func Start()  error {

	root := tview.NewTreeNode("/").
		SetColor(tcell.ColorRed)

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// 创建zk连接地址
	hosts := []string{"127.0.0.1:2181"}
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	println(conn.Server())

	list := get(conn)

	for _,v := range list {
		tmp := tview.NewTreeNode(v).
			SetColor(tcell.ColorRed)
		tmp1 := tview.NewTreeNode(v).
			SetColor(tcell.ColorRed)
		tmp.AddChild(tmp1)
		tmp.SetSelectedFunc(func() {
			tmp.SetExpanded(!tmp.IsExpanded())
		})
		root.AddChild(tmp)
	}


	//box := tview.NewBox().SetTitle("盒子").SetBorder(true)
	return tview.NewApplication().SetRoot(tree,true).EnableMouse(true).Run()
}


// 查
func get(conn *zk.Conn) []string {
	data, _, err := conn.Children("/")


	if err != nil {
		fmt.Printf("查询%s失败, err: %v\n", "/", err)
	}
	return data
}