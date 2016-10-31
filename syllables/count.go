package syllables

import "github.com/peteclark-ft/peruse/structs"

type SyllableCounter interface {
	CountSyllables(content structs.Content) int
}

type basicCounter struct{}

func NewSyllableCounter() SyllableCounter {
	return basicCounter{}
}

func (c basicCounter) CountSyllables(content structs.Content) int {
	var total float64
	for _, sentence := range content.Sentences {
		for _, word := range sentence.Words {
			total += count(word)
		}
	}

	return int(total)
}

func count(word string) float64 {
	return float64(len(word)) / 3
}
