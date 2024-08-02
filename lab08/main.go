package main

import (
	"fmt"
	"lab08/cli"
)

func main() {
	cli := cli.NewCLI()
	err := cli.Entry()
	if err != nil {
		fmt.Println(err)
		return
	}
	cli.Run()
}
