package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

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
	versionFile := flag.String("versionFile", "", "output JSON file for version pairs")
	manifestFile := flag.String("manifestFile", "", "output JSON file for manifests")
	flag.Parse()

	// Validate flags.
	if *versionFile == "" {
		fmt.Println("flag -versionFile must not be empty")
		os.Exit(1)
	}
	if *manifestFile == "" {
		fmt.Println("flag -manifestFile must not be empty")
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
	err = ioutil.WriteFile(*versionFile, data, 0644)
	if err != nil {
		panic(err)
	}

	// Generate manifests.
	rows, err = db.Query(manifestsQuery)
	if err != nil {
		panic(err)
	}

	manifests := make(map[string]map[string]string)
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
		}

		manifest := manifests[name]
		if manifest == nil {
			manifests[name] = make(map[string]string)
			manifest = manifests[name]
		}
		manifest[dependency] = spec
	}

	// Write manifests JSON file.
	manifestFd, err := os.Create(*manifestFile)
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
}
