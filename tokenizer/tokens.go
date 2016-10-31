package tokenizer

type Token int

const (
	eof           = rune(0)
	ILLEGAL Token = iota
	EOF
	FULL_STOP
	WS
	WORD
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '-'
}

func isEndOfSentence(ch rune) bool {
	return ch == '.'
}

func isText(ch rune) bool {
	return !(isWhitespace(ch) || isEndOfSentence(ch))
}

func isPunctuation(ch rune) bool {
	return ch == ',' || ch == '"' || ch == '“' || ch == '”' || ch == '’'
}
