package main

import (
	"os"
    "log"
    "fmt"
)

var cmdAdd = &Command{
	UsageLine: "add [-p HUBIC_PATH] file",
	Short:     "add file into your hubiC account",
	Long: `
add allows you to add new files into your hubiC account

`,
}

var (
    path* string
)

func init() {
	cmdAdd.Run = runAdd
    path = cmdAdd.Flag.String("p", "/", "Path to file")
}

func runAdd(cmd *Command, args []string) {
    if len(args) != 1 {
        help([]string{ cmd.Name() })
        return
    }
    filePath := args[0]

    file, err := os.Open(filePath) // For read access.
    if err != nil {
        log.Fatal(err)
    }

    fileInfo, err := os.Stat(filePath)
    if err != nil {
        log.Fatal(err)
    }

    if fileInfo.IsDir() {
        log.Fatalf("%s: not a file", fileInfo.Name())
    }

    err = Account.AddFile(*path, fileInfo.Name(), file)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%s has been successfully uploaded into your hubiC account :-)\n", fileInfo.Name())

    return
}
