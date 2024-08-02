package main

import (
	"fmt"
	"lab05/cli"
)

func main() {
	cli := cli.New()
	err := cli.Entry()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = cli.CreateFile()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("File created successfully.")
}
