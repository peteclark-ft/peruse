package syllables

import (
	"testing"

	"github.com/peteclark-ft/peruse/structs"
	"github.com/stretchr/testify/assert"
)

func TestCountSyllables(t *testing.T) {
	counter := NewSyllableCounter()

	content := structs.Content{
		Sentences: []structs.Sentence{
			{
				Words: []string{"arsenal", "are", "pretty", "kewl"},
			},
		},
	}

	total := counter.CountSyllables(content)
	assert.Equal(t, 6, total, "Expect 6")
}
