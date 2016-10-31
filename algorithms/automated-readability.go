package algorithms

import (
	"errors"

	"github.com/peteclark-ft/peruse/structs"
	"github.com/peteclark-ft/peruse/syllables"
)

type AutomatedReadabilityAnalyser interface {
	AutomatedReadability(content structs.Content) (float64, error)
}

func NewAutomatedReadabilityAnalyser(counter syllables.SyllableCounter) AutomatedReadabilityAnalyser {
	return algorithm{counter}
}

func (a algorithm) AutomatedReadability(content structs.Content) (float64, error) {
	if content.TotalSentences == 0 || content.TotalWords == 0 {
		return 0, errors.New("Failed to count syllables, or the content had none.")
	}

	return (float64(4.71) * (float64(content.TotalCharacters) / float64(content.TotalWords))) + (float64(0.5) * (float64(content.TotalWords) / float64(content.TotalSentences))) - float64(21.43), nil
}
