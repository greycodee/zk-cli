package tui

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/greycodee/zk-cli/core"
	"github.com/rivo/tview"
	"github.com/samuel/go-zookeeper/zk"
	"path/filepath"
)

type TUI struct {
	treeView *tview.TreeView
	stateView *tview.TextView
	dataView *tview.TextView
	zkClient *core.ZK
}

func NewTUI() (err error) {
	client, err := core.NewZKClient("127.0.0.1:2181")
	if err != nil {
		return
	}

	tui := TUI{
		zkClient: client,
	}
	tui.initTreeView()
	tui.initStateView()
	tui.initDataView()

	rightLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.stateView,13,0,false).
		AddItem(tui.dataView,0,3,false)

	mainLayout := tview.NewFlex().
		AddItem(tui.treeView,0,2,true).
		AddItem(rightLayout,0,3,false)
	err = tview.NewApplication().SetRoot(mainLayout, true).EnableMouse(true).Run()
	if err != nil {
		return
	}
	return
}

func (tui *TUI) treeSelectedFunc(node *tview.TreeNode)  {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}
	children := node.GetChildren()

	path := reference.(string)
	_, stat, err := tui.zkClient.GetData(path)
	if err != nil {
		panic(err)
	}

	tui.setDataViewText(stat)
	tui.dataView.SetText(fmt.Sprintf("%d",node.GetLevel()))

	if len(children) == 0 {
		tui.addTreeNode(node, path)
	} else {
		node.SetExpanded(!node.IsExpanded())

	}
}

func (tui *TUI) addTreeNode(target *tview.TreeNode, path string) {
	list,_,err := tui.zkClient.Children(path)
	if err != nil {
		panic(err)
	}
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
	rootDir := "/"
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	tree.SetSelectedFunc(tui.treeSelectedFunc)
	tree.SetTitle("节点").SetBorder(true)

	tui.addTreeNode(root,rootDir)

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

func (tui *TUI) setDataViewText(state *zk.Stat)  {
	str := bytes.Buffer{}
	str.WriteString("cZxid：")
	str.WriteString(fmt.Sprintf("%d",state.Czxid))
	str.WriteString("\n")
	str.WriteString("ctime：")
	str.WriteString(fmt.Sprintf("%d",state.Ctime))
	str.WriteString("\n")
	str.WriteString("mZxid：")
	str.WriteString(fmt.Sprintf("%d",state.Mzxid))
	str.WriteString("\n")
	str.WriteString("mtime：")
	str.WriteString(fmt.Sprintf("%d",state.Mtime))
	str.WriteString("\n")
	str.WriteString("pZxid：")
	str.WriteString(fmt.Sprintf("%d",state.Pzxid))
	str.WriteString("\n")
	str.WriteString("cversion：")
	str.WriteString(fmt.Sprintf("%d",state.Cversion))
	str.WriteString("\n")
	str.WriteString("dataVersion：")
	str.WriteString(fmt.Sprintf("%d",state.Version))
	str.WriteString("\n")
	str.WriteString("aclVersion：")
	str.WriteString(fmt.Sprintf("%d",state.Aversion))
	str.WriteString("\n")
	str.WriteString("ephemeralOwner：")
	str.WriteString(fmt.Sprintf("%d",state.EphemeralOwner))
	str.WriteString("\n")
	str.WriteString("dataLength：")
	str.WriteString(fmt.Sprintf("%d",state.DataLength))
	str.WriteString("\n")
	str.WriteString("numChildren：")
	str.WriteString(fmt.Sprintf("%d",state.NumChildren))
	str.WriteString("\n")
	tui.stateView.SetText(str.String())
}