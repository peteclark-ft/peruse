package main

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/peteclark-ft/peruse/structs"
	"github.com/urfave/cli"
)

type Report struct {
	Easiest structs.EasiestTopResults `json:"easiest,omitempty"`
	Hardest structs.HardestTopResults `json:"hardest,omitempty"`
	All     []structs.List            `json:"drilldown,omitempty"`
}

func output(ctx *cli.Context, records []structs.List) error {
	top := ctx.Int("top")
	report := Report{}

	if top > 0 {
		easiest, hardest := processTopResults(top, records)
		report.Easiest = easiest
		report.Hardest = hardest
	} else {
		report.All = records
	}

	encoder := json.NewEncoder(os.Stdout)
	return encoder.Encode(report)
}

func processTopResults(top int, records []structs.List) (structs.EasiestTopResults, structs.HardestTopResults) {
	sort.Sort(structs.ByAverage(records))
	easiest := structs.NewEasiestTopResults(top)
	hardest := structs.NewHardestTopResults(top)

	for _, record := range records {
		for _, article := range record.Content {
			easiest.Push(article)
			hardest.Push(article)
		}
	}

	return easiest, hardest
}

func getAverage(count, total float64) float64 {
	if count == 0 {
		return 0
	}
	return total / count
}
