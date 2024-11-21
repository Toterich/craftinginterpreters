package error

import (
	"fmt"
	"log"
)

func MakeErrorMsg(line int, where string, message string) error {
	return fmt.Errorf("[line %d] Error %s: %s", line, where, message)
}

func LogError(line int, message string) {
	log.Println(MakeErrorMsg(line, "", message))
}
