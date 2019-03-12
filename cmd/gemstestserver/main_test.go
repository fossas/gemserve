package main_test

import (
	"encoding/json"
	"testing"

	"github.com/Masterminds/semver"

	"github.com/fossas/gemstest/bindata"
)

type Versions struct {
	Versions [][]string
}

func TestVersionsAreValidSemver(t *testing.T) {
	fixtures, err := bindata.Asset("../../bindata/versions.json")
	if err != nil {
		t.Error(err)
	}

	var versions Versions
	err = json.Unmarshal(fixtures, &versions)
	if err != nil {
		t.Error(err)
	}

	for _, pairs := range versions.Versions {
		pkg := pairs[0]
		version := pairs[1]
		_, err := semver.NewVersion(version)
		if err != nil {
			t.Errorf("%s: %s@%s", err.Error(), pkg, version)
		}
	}
}
