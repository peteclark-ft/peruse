package structs

import (
	"sort"
	"sync"
)

type List struct {
	UUID    string    `json:"uuid"`
	Title   string    `json:"title"`
	Average float64   `json:"average"`
	Content []Article `json:"content"`
}

type ByAverage []List

func (b ByAverage) Len() int { return len(b) }

func (b ByAverage) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByAverage) Less(i, j int) bool {
	return (b[i].Average < b[j].Average)
}

type Article struct {
	WebURL string `json:"webUrl"`
	Score  Score  `json:"score"`
}

type Score struct {
	Raw                  string  `json:"raw,omitempty"`
	FleschKincaid        float64 `json:"fleschKincaid"`
	AutomatedReadability float64 `json:"automatedReadability"`
}

type ByScore []Article

func (b ByScore) Len() int { return len(b) }

func (b ByScore) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByScore) Less(i, j int) bool {
	return (b[i].Score.AutomatedReadability < b[j].Score.AutomatedReadability)
}

func NewHardestTopResults(size int) HardestTopResults {
	articles := []Article{}
	return HardestTopResults{
		Articles: articles,
		lock:     &sync.RWMutex{},
		size:     size,
	}
}

func NewEasiestTopResults(size int) EasiestTopResults {
	articles := []Article{}
	return EasiestTopResults{
		Articles: articles,
		lock:     &sync.RWMutex{},
		size:     size,
	}
}

type TopResults struct {
	Articles []Article `json:"results"`
	lock     *sync.RWMutex
	size     int
}

type EasiestTopResults TopResults

type HardestTopResults TopResults

func (t *HardestTopResults) Push(article Article) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if article.Score.AutomatedReadability < 3.4 {
		return
	}

	t.Articles = append(t.Articles, article)
	sort.Sort(ByScore(t.Articles))

	length := len(t.Articles)
	if length > t.size {
		t.Articles = t.Articles[:length-1] // remove the last entry
	}
}

func (t *EasiestTopResults) Push(article Article) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if article.Score.AutomatedReadability > 16 {
		return
	}

	t.Articles = append(t.Articles, article)
	sort.Sort(sort.Reverse(ByScore(t.Articles)))

	length := len(t.Articles)
	if length > t.size {
		t.Articles = t.Articles[:length-1] // remove the last entry
	}
}
