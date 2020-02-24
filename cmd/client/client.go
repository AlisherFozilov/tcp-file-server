package main

import (
	"fmt"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/base"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/download"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/list"
	"github.com/AlisherFozilov/file-server/cmd/client/internal/upload"
	"os"
)

const commandsList = `Here's list of gc:
upload
download
list`

func init() {
	base.Client.Commands = []*base.Command {
		download.CmdDownload,
		upload.CmdUpload,
		list.CmdList,
	}
}

func main() {
	rootCommand := base.Client
	args := os.Args[1:]

	for {
		for _, command := range rootCommand.Commands {
			if command.Name != args[0] {
				continue
			}

			if len(command.Commands) > 0 {
				rootCommand = command
				args = args[1:]
				if len(args) == 0 {
					fmt.Println(command.Name)
					fmt.Println(command.UsageLine)
					fmt.Println(command.ShortDescription)
					return
				}

				break
			}

			command.Run(command, args)
			return
		}

		fmt.Println("unknown command ", args[0])
		fmt.Println(commandsList)
		return
	}
}