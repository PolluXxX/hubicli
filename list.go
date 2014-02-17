package main

import (
	"fmt"
	"log"
)

var cmdList = &Command{
	UsageLine: "list [directory]",
	Short:     "list your files on hubiC",
	Long: `
list gives you a complete list of all your files on hubiC

`,
}

func init() {
	cmdList.Run = runList
}

func runList(cmd *Command, args []string) {
    if len(args) > 1 {
		help([]string{cmd.Name()})
		return
	}

    path := ""
    if len(args) > 0 {
        path = args[0]
    }

	files, err := Account.List(path)
	if err != nil {
		log.Fatal(err)
	}

    for _, file := range *files {
        name := file.Name
        if file.ContentType == "application/directory" {
            name += "/"
        }

        if name != "" {
            fmt.Println(name)
        }
    }
}
