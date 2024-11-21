package error

import "fmt"

func MakeErrorMsg(line int, where string, message string) error {
	return fmt.Errorf("[line %d] Error %s: %s", line, where, message)
}
