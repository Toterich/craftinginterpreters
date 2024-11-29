package assert

import "fmt"

func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

func AssertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func MissingCase(case_ any) string {
	return fmt.Sprintf("Switch misses case: %v", case_)
}
