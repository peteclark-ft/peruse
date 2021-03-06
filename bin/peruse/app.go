package main

import (
	"bufio"
	"encoding/json"
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
	app.Name = "peruse"
	app.Usage = "Readability analysis for text content"
	app.Version = "v0.0.1"

	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "json",
			Usage: "Is stdin in json content format?",
		},
	}

	app.Flags = flags

	app.Action = func(ctx *cli.Context) (err error) {
		var plaintext string
		if ctx.Bool("json") {
			plaintext, err = readJSON()
		} else {
			plaintext, err = readPlaintext()
		}

		if err != nil {
			return err
		}

		tokenizer := tokenizer.NewTokenizer(strings.NewReader(plaintext))
		content := tokenizer.Tokenize()

		counter := syllables.NewSyllableCounter()

		fleschKincaid := algorithms.NewFleschKincaidAnalyser(counter)
		automatedReadability := algorithms.NewAutomatedReadabilityAnalyser(counter)

		fk, err := fleschKincaid.FleschKincaid(content)
		ar, err := automatedReadability.AutomatedReadability(content)
		if err != nil {
			return err
		}

		score := structs.Score{
			Raw:                  content.Raw,
			FleschKincaid:        fk,
			AutomatedReadability: ar,
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.Encode(score)

		return nil
	}

	app.Run(os.Args)
}

func readJSON() (string, error) {
	decoder := json.NewDecoder(os.Stdin)

	var uppContent structs.UPPContent
	err := decoder.Decode(&uppContent)
	if err != nil {
		return "", err
	}

	stripXML, err := xml.ParseBodyXML(strings.NewReader(uppContent.BodyXML))
	if err != nil {
		return "", err
	}

	return stripXML, nil
}

func readPlaintext() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var plaintext string
	for scanner.Scan() {
		plaintext = plaintext + scanner.Text()
	}
	return plaintext, nil
}
