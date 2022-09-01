package err

import "errors"

var (
	ErrUnknownArgs = errors.New("unknown arguments")
	ErrArgsConflict = errors.New("arguments conflict")
	ErrArgsLack = errors.New("lack of arguments")
	ErrIllFormedNS = errors.New("ill formed argument 'ns'") 
	ErrPortOutRange = errors.New("port is out of range")
	ErrIllFormedIP = errors.New("ill formed ip")
	ErrUnknownCmd = errors.New("unknown command")
	ErrLevelDBInit = errors.New("initiate leveldb failed")
	ErrRunPython = errors.New("run python failed")
)