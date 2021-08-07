package main

import (
	"flag"
	"github.com/greycodee/zk-cli/tui"
)

var host = flag.String("h","127.0.0.1:2181","zookeeper server address")

func init() {
	flag.Parse()
}

func main()  {
	tui.NewTUI().Run(*host)
}
