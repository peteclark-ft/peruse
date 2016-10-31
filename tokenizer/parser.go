package tokenizer

import (
	"bufio"
	"io"

	"github.com/peteclark-ft/peruse/structs"
)

type tokenizer struct {
	lexer *lexer
}

// ContentTokenizer tokenizes a text document into words sentences and totals
type ContentTokenizer interface {
	Tokenize() structs.Content
}

// NewTokenizer returns a new tokenizer
func NewTokenizer(reader io.Reader) ContentTokenizer {
	lexer := &lexer{bufio.NewReader(reader)}
	return tokenizer{lexer: lexer}
}

func (t tokenizer) Tokenize() structs.Content {
	content := [][]string{{}}

	for {
		token, value := t.lexer.Scan()
		if token == EOF {
			break
		}

		if token == WORD {
			currentSentence := content[len(content)-1]
			content[len(content)-1] = append(currentSentence, value)
		}

		if token == WS {
			continue
		}

		if token == FULL_STOP {
			content = append(content, []string{})
		}
	}

	var sentences []structs.Sentence
	var totalWords int
	var totalCharacters int

	for _, s := range content {
		sentences = append(sentences, structs.Sentence{Words: s})

		for _, w := range s {
			totalCharacters += len(w)
		}

		totalWords += len(s)
	}

	return structs.Content{
		TotalWords:      totalWords,
		TotalSentences:  len(sentences),
		TotalCharacters: totalCharacters,
		Sentences:       sentences,
	}
}
