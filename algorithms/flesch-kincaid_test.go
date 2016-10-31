package algorithms

import (
	"testing"

	"github.com/peteclark-ft/peruse/structs"
	"github.com/peteclark-ft/peruse/syllables"
	"github.com/stretchr/testify/assert"
)

func TestFleschKincaid(t *testing.T) {
	analyser := NewFleschKincaidAnalyser(syllables.NewSyllableCounter())
	content := structs.Content{
		TotalSentences: 1,
		TotalWords:     4,
		Sentences: []structs.Sentence{
			{
				Words: []string{"arsenal", "are", "pretty", "kewl"},
			},
		},
	}

	verdict, err := analyser.FleschKincaid(content)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 75, int(verdict), "Should be 76")
}

func TestFails(t *testing.T) {
	analyser := NewFleschKincaidAnalyser(syllables.NewSyllableCounter())
	content := structs.Content{
		TotalSentences: 0,
		TotalWords:     4,
		Sentences: []structs.Sentence{
			{
				Words: []string{"arsenal", "are", "pretty", "kewl"},
			},
		},
	}

	_, err := analyser.FleschKincaid(content)
	if err == nil {
		t.Fatal(err)
	}
}
