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
	ErrIllFormedTime = errors.New("ill formed time")
)

func ErrMatch(e error) {
	switch e {
		case ErrArgsConflict:
			ErrHandle()
		case ErrArgsLack:
			ErrHandle()
		case ErrIllFormedIP:
			ErrHandle()
		case ErrIllFormedNS:
			ErrHandle()
		case ErrLevelDBInit:
			ErrHandle()
		case ErrUnknownArgs:
			ErrHandle()
		case ErrPortOutRange:
			ErrHandle()
		default:
	}
}



func ErrHandle()  {

}