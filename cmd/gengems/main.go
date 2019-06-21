package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const versionsQuery = `
SELECT
	r.name AS package,
	v.number AS version
FROM rubygems r
JOIN versions v ON r.id = v.rubygem_id;
`

const manifestsQuery = `
SELECT
  r.name AS package,
  v.number AS version,
  rd.name AS dependency,
  d.requirements AS spec,
  d.scope
FROM versions v
JOIN rubygems r ON v.rubygem_id = r.id
JOIN dependencies d ON d.version_id = v.id
JOIN rubygems rd ON rd.id = d.rubygem_id;
`

func main() {
	// Set help text and flags.
	flag.Usage = func() {
		fmt.Printf(`%s generates data files for gemstestserver.

Usage:

`, os.Args[0])
		flag.PrintDefaults()
	}
	pg := flag.String("pg", "postgresql://postgres@localhost:5432/rubygems?sslmode=disable", "postgres URL for RubyGems dump")
	dir := flag.String("dir", "", "output directory")
	flag.Parse()

	// Validate flags.
	if *dir == "" {
		fmt.Println("flag -versionFile must not be empty")
		os.Exit(1)
	}

	// Connect to dump database.
	db := sqlx.MustConnect("postgres", *pg)

	// Generate version pairs.
	rows, err := db.Query(versionsQuery)
	if err != nil {
		panic(err)
	}

	var versions [][2]string
	for rows.Next() {
		var name, version string
		err := rows.Scan(&name, &version)
		if err != nil {
			panic(err)
		}

		// Ignore non-semver versions (they mess up `/compare` and `/within`).
		_, err = semver.NewVersion(version)
		if err != nil {
			fmt.Printf("Ignoring bad version: %s %s\n", name, version)
			continue
		}

		versions = append(versions, [2]string{name, version})
	}

	// Write versions JSON file.
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filepath.Join(*dir, "versions.json"), data, 0644)
	if err != nil {
		panic(err)
	}

	// Generate manifests.
	rows, err = db.Query(manifestsQuery)
	if err != nil {
		panic(err)
	}

	manifests := make(map[string]map[string]map[string]string)
	for rows.Next() {
		var name, version, dependency, spec, scope string
		err := rows.Scan(&name, &version, &dependency, &spec, &scope)
		if err != nil {
			panic(err)
		}

		// Ignore non-semver versions and specs.
		_, err = semver.NewVersion(version)
		if err != nil {
			fmt.Printf("Ignoring bad version: %s %s\n", name, version)
			continue
		}
		_, err = semver.NewConstraint(spec)
		if err != nil {
			fmt.Printf("Ignoring bad spec: %s %s depends on %s %s (%s)\n", name, version, dependency, spec, scope)
			continue
		}

		pkg := manifests[name]
		if pkg == nil {
			manifests[name] = make(map[string]map[string]string)
			pkg = manifests[name]
		}
		deps := pkg[version]
		if deps == nil {
			pkg[version] = make(map[string]string)
			deps = pkg[version]
		}
		deps[dependency] = spec
	}

	// Write manifests JSON file.
	manifestFd, err := os.Create(filepath.Join(*dir, "manifests.json"))
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(manifestFd)
	enc.SetEscapeHTML(false) // See https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
	enc.SetIndent("", "  ")
	err = enc.Encode(manifests)
	if err != nil {
		panic(err)
	}

	// Generate manifest snapshot using the dependencies of the latest version for
	// each package. This snapshot is NOT guaranteed to have a resolvable graph
	// for any packages.
	snapshot := make(map[string]map[string]string)
	for pkg, versions := range manifests {
		latestVersion := &semver.Version{}
		latestDeps := make(map[string]string)
		for version, deps := range versions {
			v, err := semver.NewVersion(version)
			if err != nil {
				panic(err)
			}
			if v.GreaterThan(latestVersion) {
				latestVersion = v
				latestDeps = deps
			}
		}
		snapshot[pkg] = latestDeps
	}

	// Write snapshots JSON file.
	snapshotFd, err := os.Create(filepath.Join(*dir, "snapshot.json"))
	if err != nil {
		panic(err)
	}
	enc = json.NewEncoder(snapshotFd)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(snapshot)
	if err != nil {
		panic(err)
	}
}
