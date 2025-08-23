package core

import (
	"errors"
	"strconv"
	"syscall"
	"time"

	"github.com/lyxuansang91/redis-crash-course/internal/constant"
	"github.com/lyxuansang91/redis-crash-course/internal/data_structure"
)

type CommandExecutor interface {
	Ping(args []string) []byte
	Set(args []string) []byte
	Get(args []string) []byte
	Ttl(args []string) []byte
	Expire(args []string) []byte
	Del(args []string) []byte
	Exists(args []string) []byte
	ExecuteAndResponse(command *Command, connFd int) error
}

type CommandExecutorImpl struct {
	dictStore *data_structure.Dict
}

func NewCommandExecutor(dictStore *data_structure.Dict) CommandExecutor {
	return &CommandExecutorImpl{
		dictStore: dictStore,
	}
}

func (cmd *CommandExecutorImpl) Ping(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	return res
}

func (cmd *CommandExecutorImpl) Set(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'SET' command"), false)
	}

	var key, value string
	var ttlMs int64 = -1

	key, value = args[0], args[1]
	if len(args) > 2 {
		ttlTime, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		if ttlTime < 0 {
			return Encode(errors.New("(error) ERR invalid TTL value"), false)
		}
		switch args[2] {
		case "PX":
			ttlMs = ttlTime
		case "EX":
			ttlMs = ttlTime * 1000
		case "EXAT":
			ttlMs = ttlTime*1000 - time.Now().UnixMilli()
		case "PXAT":
			ttlMs = ttlTime - time.Now().UnixMilli()
		default:
			return Encode(errors.New("(error) ERR invalid TTL unit"), false)
		}

		if ttlMs <= 0 {
			return Encode(errors.New("(error) ERR invalid TTL value"), false)
		}

	}

	cmd.dictStore.Set(key, cmd.dictStore.NewObj(key, value, ttlMs))
	return constant.RespOk
}

func (cmd *CommandExecutorImpl) Get(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}

	key := args[0]
	obj := cmd.dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}

	return Encode(obj.Value, false)
}

func (cmd *CommandExecutorImpl) Ttl(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'TTL' command"), false)
	}
	key := args[0]
	obj := cmd.dictStore.Get(key)
	if obj == nil {
		return constant.TtlKeyNotExist
	}

	exp, isExpirySet := cmd.dictStore.GetExpiry(key)
	if !isExpirySet {
		return constant.TtlKeyExistNoExpire
	}

	remainMs := exp - time.Now().UnixMilli()
	if remainMs < 0 {
		return constant.TtlKeyNotExist
	}

	return Encode(int64(remainMs/1000), false)
}

func (cmd *CommandExecutorImpl) Expire(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'EXPIRE' command"), false)
	}

	key := args[0]
	ttlSeconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
	}

	if ttlSeconds < 0 {
		return Encode(errors.New("(error) ERR invalid TTL value"), false)
	}

	if cmd.dictStore.Get(key) == nil {
		return constant.ResIntegerNotOk
	}

	ttlMs := ttlSeconds * 1000

	cmd.dictStore.SetExpiry(key, ttlMs)
	return constant.ResIntegerOk
}

func (cmd *CommandExecutorImpl) ExpireAt(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'EXPIREAT' command"), false)
	}
	return []byte("-CMD NOT FOUND\r\n")
}

// ExecuteAndResponse given a Command, executes it and responses
func (cmd *CommandExecutorImpl) ExecuteAndResponse(command *Command, connFd int) error {
	var res []byte

	switch command.Cmd {
	case CmdPing:
		res = cmd.Ping(command.Args)
	case CmdSet:
		res = cmd.Set(command.Args)
	case CmdGet:
		res = cmd.Get(command.Args)
	case CmdTtl:
		res = cmd.Ttl(command.Args)
	case CmdExpire:
		res = cmd.Expire(command.Args)
	case CmdExpireAt:
		res = cmd.ExpireAt(command.Args)
	case CmdDel:
		res = cmd.Del(command.Args)
	case CmdExists:
		res = cmd.Exists(command.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}
	_, err := syscall.Write(connFd, res)
	return err
}

func (cmd *CommandExecutorImpl) Exists(args []string) []byte {
	count := 0
	for _, key := range args {
		if cmd.dictStore.Get(key) != nil  {
			count++
		}
	}
	return Encode(int64(count), false)
}

func (cmd *CommandExecutorImpl) Del(args []string) []byte {
	count := 0
	for _, key := range args {
		if cmd.dictStore.Get(key) != nil && cmd.dictStore.Del(key) {
			count++
		}
	}
	return Encode(int64(count), false)
}
