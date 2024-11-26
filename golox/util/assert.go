package util

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
