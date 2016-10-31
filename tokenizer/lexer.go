package tokenizer

import (
	"bufio"
	"bytes"
)

type lexer struct {
	stream *bufio.Reader
}

func (l *lexer) read() rune {
	ch, _, err := l.stream.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (l *lexer) unread() {
	l.stream.UnreadRune()
}

func (l *lexer) Scan() (Token, string) {
	ch := l.read()

	if ch == eof {
		return EOF, string(ch)
	}

	if isWhitespace(ch) {
		l.unread()
		return WS, l.readType(isWhitespace, isPunctuation)
	}

	if isEndOfSentence(ch) {
		l.unread()
		result := l.readType(isEndOfSentence, isPunctuation)
		if len(result) == 3 {
			return WS, ""
		}

		return FULL_STOP, string(ch)
	}

	if isText(ch) {
		l.unread()
		return WORD, l.readType(isText, isPunctuation)
	}

	return ILLEGAL, string(ch)
}

func (l *lexer) readType(check func(ch rune) bool, ignore func(ch rune) bool) string {
	var buf bytes.Buffer
	buf.WriteRune(l.read())

	for {
		ch := l.read()

		if ch == eof {
			break
		}

		if ignore(ch) {
			continue
		}

		if !check(ch) {
			l.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return buf.String()
}
