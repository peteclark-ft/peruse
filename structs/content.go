package structs

type Content struct {
	TotalWords      int
	TotalSentences  int
	TotalCharacters int
	Sentences       []Sentence
}

type Sentence struct {
	Words []string
}

type UPPContent struct {
	BodyXML string `json:"bodyXml"`
	ID      string `json:"id"`
	WebURL  string `json:"webUrl"`
	Type    string `json:"type"`
}
