package algorithms

import (
	"errors"

	"github.com/peteclark-ft/peruse/structs"
	"github.com/peteclark-ft/peruse/syllables"
)

type algorithm struct {
	counter syllables.SyllableCounter
}

type FleschKincaidAnalyser interface {
	FleschKincaid(content structs.Content) (float64, error)
}

func NewFleschKincaidAnalyser(counter syllables.SyllableCounter) FleschKincaidAnalyser {
	return algorithm{counter}
}

func (a algorithm) FleschKincaid(content structs.Content) (float64, error) {
	totalSyllables := a.counter.CountSyllables(content)

	if totalSyllables == 0 || content.TotalSentences == 0 || content.TotalWords == 0 {
		return 0, errors.New("Failed to count syllables, or the content had none.")
	}

	return float64(206.835) - (float64(1.015) * (float64(content.TotalWords) / float64(content.TotalSentences))) - (float64(84.6) * (float64(totalSyllables) / float64(content.TotalWords))), nil
}
