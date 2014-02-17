package main

import (
	"fmt"
	"log"
)

var cmdList = &Command{
	UsageLine: "list",
	Short:     "list your files on hubiC",
	Long: `
list gives you a complete list of all your files on hubiC

`,
}

func init() {
	cmdList.Run = runList
}

func runList(cmd *Command, args []string) {
	files, err := Account.List("/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", files)
}
