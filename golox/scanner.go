package main

type Scanner struct {
	source string
}

func (s Scanner) scanTokens() []string {
	return []string{s.source}
}
