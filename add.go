package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KDevroede/monio"
	"github.com/cheggaaa/pb"
)

var cmdAdd = &Command{
	UsageLine: "add [-p HUBIC_PATH] [-n FILENAME] file",
	Short:     "add file into your hubiC account",
	Long: `
add allows you to add new files into your hubiC account

The flags are:
        -p
            Choose path where your file will be stored.
            Must start and end with '/' character.
            If path does not exists, it will be created first.

        -n
            By default, your file will be uploaded with its own name.
            You can rename it by specify this flag, followed by new 
            file name.
            This option can not be used if multiple files are specified.


Files can be piped into hubicli command :
    grep 'HelloWorld' file.txt | hubicli add -n Hello.txt

If -n is not specified when used with pipe, the file will be 
uploaded with Unix timestamp as file name

`,
}

var (
	pathA     *string
	fileNameA *string
)

func init() {
	cmdAdd.Run = runAdd
	pathA = cmdAdd.Flag.String("p", "", "Path to file")
	fileNameA = cmdAdd.Flag.String("n", "", "Filename")
}

func runAdd(cmd *Command, args []string) {
	fileName := *fileNameA
	path := *pathA

	if len(args) == 0 {
		if fileName == "" {
			fileName = fmt.Sprintf("%d", time.Now().Unix())
		}
		args = append(args, "/dev/stdin")
	}

	if path == "" {
		path = "/"
	}

	if len(args) != 1 && fileName != "" {
		log.Fatal("You can not specify -n option with multiple files")
	}

	if path[len(path)-1:len(path)] != "/" || path[0:1] != "/" {
		log.Fatalf("Path must start and end with '/' char. Found %s", path)
	}

	for _, filePath := range args {
		file, err := os.Open(filePath)
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

		c := make(chan int)
		reader := monio.NewReader(file, c)

		name := fileName
		if name == "" {
			name = fileInfo.Name()
		}

		fmt.Printf("Start uploading %s\n", name)
		go func() {
			err = Account.AddFile(path, name, reader)
			if err != nil {
				log.Fatal(err)
			}
		}()

		bar := pb.StartNew(int(fileInfo.Size()))
		bar.SetUnits(pb.U_BYTES)
		for bytes := range c {
			bar.Add(bytes)
		}
		bar.FinishPrint(fmt.Sprintf("%s uploaded!", name))
	}

	return
}
