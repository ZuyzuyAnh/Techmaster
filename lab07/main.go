package main

import (
	"fmt"
	"lab07/cli"
)

func main() {
	cli := cli.NewCLI()
	err := cli.Entry()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	cli.SearchInDirectory()
}
