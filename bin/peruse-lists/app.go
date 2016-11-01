package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/peteclark-ft/peruse/algorithms"
	"github.com/peteclark-ft/peruse/structs"
	"github.com/peteclark-ft/peruse/syllables"
	"github.com/peteclark-ft/peruse/tokenizer"
	"github.com/peteclark-ft/peruse/xml"
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
		cli.IntFlag{
			Name:  "top",
			Usage: "Show the top (and bottom) x scores, and the lists they're in.",
			Value: -1,
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

		lists := []structs.List{}
		for _, uuid := range uuids {
			list, items, _ := get.requestList(uuid.UUID)

			var articles []structs.Article

			var total float64
			for _, content := range items {
				stripXML, err := xml.ParseBodyXML(strings.NewReader(content.BodyXML))
				if err != nil {
					log.Println(err)
					continue
				}

				tokenizer := tokenizer.NewTokenizer(strings.NewReader(stripXML))
				tokenized := tokenizer.Tokenize()

				fk, _ := fleschKincaid.FleschKincaid(tokenized)
				ar, _ := automatedReadability.AutomatedReadability(tokenized)

				total += ar

				article := structs.Article{
					WebURL: content.WebURL,
					Score: structs.Score{
						FleschKincaid:        fk,
						AutomatedReadability: ar,
					},
				}
				articles = append(articles, article)
			}

			average := getAverage(float64(len(articles)), total)

			lists = append(lists, structs.List{
				UUID:    uuid.UUID,
				Average: average,
				Title:   list.Title,
				Content: articles,
			})
		}

		output(ctx, lists)
		return err
	}

	app.Run(os.Args)
}
