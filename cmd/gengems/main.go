package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/Masterminds/semver"
)

type Versions struct {
	Versions [][]string `json:"versions"`
}

func main() {
	infile := flag.String("input", "", "the input version JSON file")
	outfile := flag.String("output", "", "the cleaned version JSON file")
	flag.Parse()

	fixtures, err := ioutil.ReadFile(*infile)
	if err != nil {
		panic(err)
	}

	var versions Versions
	err = json.Unmarshal(fixtures, &versions)
	if err != nil {
		panic(err)
	}

	var cleaned Versions
	for _, pairs := range versions.Versions {
		version := pairs[1]
		_, err := semver.NewVersion(version)
		if err != nil {
			continue
		} else {
			cleaned.Versions = append(cleaned.Versions, pairs)
		}
	}

	data, err := json.MarshalIndent(cleaned, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(*outfile, data, 0644)
	if err != nil {
		panic(err)
	}
}
