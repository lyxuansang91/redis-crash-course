package constant

import "time"

var RespNil = []byte("$-1\r\n")
var RespExpired = []byte("$-2\r\n")
var ResIntegerOk = []byte(":1\r\n")
var ResIntegerNotOk = []byte(":0\r\n")
var RespOk = []byte("+OK\r\n")
var TtlKeyNotExist = []byte(":-2\r\n")
var TtlKeyExistNoExpire = []byte(":-1\r\n")
var ActiveExpireFrequency = 100 * time.Millisecond
var ActiveExpireSampleSize = 20
var ActiveExpireThreshold = 0.1