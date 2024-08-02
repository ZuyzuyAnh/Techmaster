package main

import (
	"fmt"
	"lab06/cli"
)

func main() {
	cli := cli.NewCLI()
	err := cli.Entry()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	cli.Run()
}
