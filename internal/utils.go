package internal

import "runtime"

func HandleError(message string, err error) {
	buf := make([]byte, 1024)
	runtime.Stack(buf, false)
	panic(message + "\n" + err.Error() + "\n" + string(buf))
}
