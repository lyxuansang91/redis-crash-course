package core

type Command struct {
	Cmd string 
	Args []string
}

const (
	CmdPing = "PING"
	CmdSet = "SET"
	CmdGet = "GET"
	CmdTtl = "TTL"
)