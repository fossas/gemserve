package main_test

import (
	"encoding/json"
	"testing"

	"github.com/Masterminds/semver"

	"github.com/fossas/gemstest/bindata"
)

func TestVersionsAreValidSemver(t *testing.T) {
	fixtures, err := bindata.Asset("../../bindata/data/versions.json")
	if err != nil {
		t.Error(err)
	}

	var versions [][]string
	err = json.Unmarshal(fixtures, &versions)
	if err != nil {
		t.Error(err)
	}

	for _, pairs := range versions {
		pkg := pairs[0]
		version := pairs[1]
		_, err := semver.NewVersion(version)
		if err != nil {
			t.Errorf("%s: %s@%s", err.Error(), pkg, version)
		}
	}
}

func TestSnapshotUsesValidSpecs(t *testing.T) {
	fixtures, err := bindata.Asset("../../bindata/data/snapshot.json")
	if err != nil {
		t.Error(err)
	}

	var snapshot map[string]map[string]string
	err = json.Unmarshal(fixtures, &snapshot)
	if err != nil {
		t.Error(err)
	}

	for pkg, specs := range snapshot {
		for dep, spec := range specs {
			_, err := semver.NewConstraint(spec)
			if err != nil {
				t.Errorf("%s: %s depending on %s@%s", err.Error(), pkg, dep, spec)
			}
		}
	}
}
