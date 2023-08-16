package cmd

import (
	Get "github.com/connorbaker/godfs/cmd/get"
	Init "github.com/connorbaker/godfs/cmd/init"
	List "github.com/connorbaker/godfs/cmd/list"
	Put "github.com/connorbaker/godfs/cmd/put"
	"github.com/urfave/cli/v2"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		Get.GetCommand(),
		Init.InitCommand(),
		List.ListCommand(),
		Put.PutCommand(),
	}
}
