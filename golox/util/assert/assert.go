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
	// It would be great if we could panic here directly, instead of passing the assertion string
	// and then calling panic() at the call site. But then golang complains about missing return types
	return fmt.Sprintf("Switch misses case: %v", case_)
}
