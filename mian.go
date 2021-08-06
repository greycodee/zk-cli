package main

import "github.com/greycodee/zk-cli/tui"

func main()  {
	//fmt.Println("hello world")

	err := tui.NewTUI()
	if err != nil {
		panic(err)
	}
}
