package base

type Command struct {
	Name             string
	Run              func(cmd *Command, args []string)
	UsageLine        string
	ShortDescription string
	Commands         []*Command
}

var Client = &Command{
	Run:              nil,
	UsageLine:        "client",
	ShortDescription: "client is a tool for keeping files on server",
}
