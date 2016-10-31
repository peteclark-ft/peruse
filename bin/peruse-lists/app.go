package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/peteclark-ft/peruse/algorithms"
	"github.com/peteclark-ft/peruse/syllables"
	"github.com/peteclark-ft/peruse/tokenizer"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "peruse-lists"
	app.Usage = "Analyses readability for all articles in a list"
	app.Version = "v0.0.1"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "api-key",
			Usage: "API Key for test.api.ft.com",
		},
		cli.StringFlag{
			Name:  "url",
			Usage: "Content environment url",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Basic auth user",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "Basic auth password",
		},
	}
	app.Flags = flags

	app.Action = func(ctx *cli.Context) error {
		decoder := json.NewDecoder(os.Stdin)
		var uuids = []struct {
			UUID string `json:"uuid"`
		}{}
		err := decoder.Decode(&uuids)
		if err != nil {
			return err
		}

		client := &http.Client{}
		get := getList{
			client:   client,
			apiKey:   ctx.String("api-key"),
			user:     ctx.String("user"),
			password: ctx.String("password"),
			url:      ctx.String("url"),
		}

		counter := syllables.NewSyllableCounter()
		fleschKincaid := algorithms.NewFleschKincaidAnalyser(counter)
		automatedReadability := algorithms.NewAutomatedReadabilityAnalyser(counter)

		var easiest *ListItem
		var hardest *ListItem

		output := []Output{}
		for _, uuid := range uuids {
			list, items, _ := get.requestList(uuid.UUID)

			var jsonItems []ListItem

			var total float64
			for _, content := range items {
				tokenizer := tokenizer.NewTokenizer(strings.NewReader(content.BodyXML))
				tokenized := tokenizer.Tokenize()

				fk, _ := fleschKincaid.FleschKincaid(tokenized)
				ar, _ := automatedReadability.AutomatedReadability(tokenized)

				total += ar

				listItem := ListItem{
					WebURL: content.WebURL,
					Score: Score{
						FleschKincaid:        fk,
						AutomatedReadability: ar,
					},
				}
				jsonItems = append(jsonItems, listItem)

				if easiest == nil || easiest.Score.AutomatedReadability > listItem.Score.AutomatedReadability {
					easiest = &listItem
				}

				if hardest == nil || hardest.Score.AutomatedReadability > listItem.Score.AutomatedReadability {
					hardest = &listItem
				}
			}

			average := getAverage(float64(len(jsonItems)), total)

			output = append(output, Output{
				UUID:    uuid.UUID,
				Average: average,
				Title:   list.Title,
				Content: jsonItems,
			})
		}

		encoder := json.NewEncoder(os.Stdout)
		sort.Sort(ByAverage(output))

		err = encoder.Encode(Report{
			Min: *easiest,
			Max: *hardest,
			All: output,
		})
		return err
	}

	app.Run(os.Args)
}

func getAverage(count, total float64) float64 {
	if count == 0 {
		return 0
	}
	return total / count
}

type Report struct {
	Min ListItem `json:"easiest"`
	Max ListItem `json:"hardest"`
	All []Output `json:"drilldown"`
}

type Output struct {
	UUID    string     `json:"uuid"`
	Title   string     `json:"title"`
	Average float64    `json:"average"`
	Content []ListItem `json:"items"`
}

type ListItem struct {
	WebURL string `json:"webUrl"`
	Score  Score  `json:"score"`
}

type Score struct {
	FleschKincaid        float64 `json:"fleschKincaid"`
	AutomatedReadability float64 `json:"automatedReadability"`
}

type ByAverage []Output

func (b ByAverage) Len() int { return len(b) }

func (b ByAverage) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByAverage) Less(i, j int) bool {
	return (b[i].Average < b[j].Average)
}
