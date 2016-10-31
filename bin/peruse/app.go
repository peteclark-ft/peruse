package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/peteclark-ft/peruse/algorithms"
	"github.com/peteclark-ft/peruse/structs"
	"github.com/peteclark-ft/peruse/syllables"
	"github.com/peteclark-ft/peruse/tokenizer"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "peruse"
	app.Usage = "Readability analysis for text content"
	app.Version = "v0.0.1"

	/*flags := []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "./config.yml",
			Usage: "Path to the YAML config file.",
		},
	}*/
	//app.Flags = flags

	app.Action = func(ctx *cli.Context) error {
		decoder := json.NewDecoder(os.Stdin)

		var uppContent structs.UPPContent
		err := decoder.Decode(&uppContent)
		if err != nil {
			return err
		}

		tokenizer := tokenizer.NewTokenizer(strings.NewReader(uppContent.BodyXML))
		content := tokenizer.Tokenize()

		counter := syllables.NewSyllableCounter()

		fleschKincaid := algorithms.NewFleschKincaidAnalyser(counter)
		automatedReadability := algorithms.NewAutomatedReadabilityAnalyser(counter)

		fk, err := fleschKincaid.FleschKincaid(content)
		ar, err := automatedReadability.AutomatedReadability(content)
		if err != nil {
			return err
		}

		score := score{
			FleschKincaid:        fk,
			AutomatedReadability: ar,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.Encode(score)

		return nil
	}

	app.Run(os.Args)
}

type score struct {
	FleschKincaid        float64 `json:"fleschKincaid"`
	AutomatedReadability float64 `json:"automatedReadability"`
}
