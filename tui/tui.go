package tui

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/greycodee/zk-cli/core"
	"github.com/greycodee/zk-cli/log"
	"github.com/rivo/tview"
	"github.com/samuel/go-zookeeper/zk"
	"net/url"
	"path/filepath"
)

type TUI struct {
	treeView *tview.TreeView
	treeRootNode *tview.TreeNode
	stateView *tview.TextView
	dataView *tview.TextView
	reFlushBtn *tview.Button

	tuiLog *log.Log
	zkClient *core.ZK

	currPath string
	zkHost string

}

const root = "/"

func NewTUI() *TUI{
	tui := TUI{}

	tui.initTreeView()
	tui.initStateView()
	tui.initDataView()
	tui.initReFlushBtn()

	return &tui
}

func (tui *TUI) Run(zkHost string)  {
	tui.tuiLog = log.NewLogs()
	client := core.NewZKClient(zkHost,tui.tuiLog)

	tui.zkClient = client
	tui.zkHost = zkHost

	tui.addTreeNode(tui.treeRootNode,root)

	leftLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.reFlushBtn,3,0,false).
		AddItem(tui.treeView,0,2,true)

	rightLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.stateView,0,1,false).
		AddItem(tui.dataView,0,1,false)

	mainView := tview.NewFlex().
		AddItem(leftLayout,0,2,true).
		AddItem(rightLayout,0,3,false)

	err := tview.NewApplication().SetRoot(mainView, true).EnableMouse(true).Run()
	if err != nil {
		panic(err)
	}

}

func (tui *TUI) treeSelectedFunc(node *tview.TreeNode)  {
	reference := node.GetReference()
	if reference == nil {
		return
	}
	children := node.GetChildren()

	path := reference.(string)
	unescape, err := url.QueryUnescape(path)
	if err != nil {
		return
	}

	tui.currPath = unescape
	data, stat := tui.zkClient.GetData(path)

	tui.stateView.Clear()
	tui.dataView.Clear()

	tui.setDataViewText(stat)
	tui.dataView.SetText(data)

	if len(children) == 0 {
		tui.addTreeNode(node, path)
	}else {
		node.SetExpanded(!node.IsExpanded())
	}
}

func (tui *TUI) addTreeNode(target *tview.TreeNode, path string) {
	list,_ := tui.zkClient.Children(path)

	for _, n := range list {

		flag := tui.zkClient.HaveChildren(filepath.Join(path, n))
		node := tview.NewTreeNode(n).
			SetReference(filepath.Join(path, n))
		if flag {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}
}

func (tui *TUI) initTreeView()  {
	tui.treeRootNode = tview.NewTreeNode(root).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(tui.treeRootNode).
		SetCurrentNode(tui.treeRootNode)
	tree.SetSelectedFunc(tui.treeSelectedFunc)
	tree.SetTitle("节点").SetBorder(true)

	tui.treeView = tree
}

func (tui *TUI) initStateView()  {
	state := tview.NewTextView()
	state.SetTitle("状态").SetBorder(true)

	tui.stateView = state
}

func (tui *TUI) initDataView()  {
	dataView := tview.NewTextView()
	dataView.SetTitle("数据").SetBorder(true)

	tui.dataView = dataView
}

func (tui *TUI) initReFlushBtn()  {
	btn := tview.NewButton("刷新当前节点")
	btn.SetBorder(true)
	btn.SetSelectedFunc(tui.flush)

	tui.reFlushBtn = btn
}

func (tui *TUI) setDataViewText(state *zk.Stat)  {
	str := bytes.Buffer{}
	str.WriteString(fmt.Sprintf("cZxid：%d\n",state.Czxid))
	str.WriteString(fmt.Sprintf("ctime：%d\n",state.Ctime))
	str.WriteString(fmt.Sprintf("mZxid：%d\n",state.Mzxid))
	str.WriteString(fmt.Sprintf("mtime：%d\n",state.Mtime))
	str.WriteString(fmt.Sprintf("pZxid：%d\n",state.Pzxid))
	str.WriteString(fmt.Sprintf("cversion：%d\n",state.Cversion))
	str.WriteString(fmt.Sprintf("dataVersion：%d\n",state.Version))
	str.WriteString(fmt.Sprintf("aclVersion：%d\n",state.Aversion))
	str.WriteString(fmt.Sprintf("ephemeralOwner：%d\n",state.EphemeralOwner))
	str.WriteString(fmt.Sprintf("dataLength：%d\n",state.DataLength))
	str.WriteString(fmt.Sprintf("numChildren：%d\n",state.NumChildren))

	str.WriteString(fmt.Sprintf("zk服务地址：%s\n",tui.zkHost))
	str.WriteString(fmt.Sprintf("当前路径：%s\n",tui.currPath))
	tui.stateView.SetText(str.String())
}

func (tui *TUI) flush(){
	currNode := tui.treeView.GetCurrentNode()
	if currNode.GetText() != root {
		reference := tui.treeView.GetCurrentNode().GetReference()
		path := reference.(string)
		currNode.ClearChildren()
		tui.addTreeNode(currNode,path)

		data, stat := tui.zkClient.GetData(path)
		tui.stateView.Clear()
		tui.dataView.Clear()
		tui.setDataViewText(stat)
		tui.dataView.SetText(data)
	}
}