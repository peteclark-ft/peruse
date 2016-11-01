package xml

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// ParseBodyXML parses content body xml into a single contiguous article string
func ParseBodyXML(reader io.Reader) (string, error) {
	content := ""

	z := html.NewTokenizer(reader)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}

		if tt == html.EndTagToken && z.Token().Data == "p" {
			if !strings.HasSuffix(strings.TrimSpace(content), ".") {
				content += "."
			}
			continue
		}

		if tt != html.TextToken {
			continue
		}

		token := z.Token()

		content += token.Data
	}

	return content, nil
}
