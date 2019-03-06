package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/Masterminds/semver"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/fossas/gemstest/bindata"
)

// CompareRequest is a request for the /compare endpoint. This endpoint takes
// two versions, a and b, and compares them.
type CompareRequest struct {
	A string
	B string
}

// WithinRequest is a request for the /within endpoint. This endpoint takes a
// version and a spec and checks whether the version is within the spec.
type WithinRequest struct {
	Version string
	Spec    string
}

//go:generate go-bindata -modtime 0 -o ../../bindata/bindata.go -pkg bindata ../../bindata/data

func main() {
	flag.Usage = func() {
		fmt.Printf(`%s is the test server for the Gems challenge.

This HTTP server listens on three paths:

  GET /versions

    Returns a JSON object whose key "versions" contains a list of (package name,
    package version) pairs.

    Example output:

      {
        "versions": [
          ["name", "1.0.0"],
          ["name", "1.0.1"],
          ["another name", "2.3.4"]
        ]
      }

  POST /compare

    Expects a JSON request body with two strings "a" and "b".
    Returns -1 if a < b, 0 if a == b, and 1 if a > b.

    Example input:

      {
        "a": "1.0.0",
        "b": "1.0.1"
      }

    Example output:

      -1

  POST /within

    Expects a JSON request body with two strings "version" and "spec". Returns
    true if the version is within the spec, and false otherwise.

    Example input:

      {
        "version": "1.0.2",
        "spec": "^1.0.0"
      }

    Example output:

      true

Flags:

`, os.Args[0])
		flag.PrintDefaults()
	}

	port := flag.Int("port", 8000, "the port to listen on")
	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/compare", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var req CompareRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			panic(err)
		}

		a := semver.MustParse(req.A)
		b := semver.MustParse(req.B)

		_, err = w.Write([]byte(strconv.Itoa(a.Compare(b))))
		if err != nil {
			panic(err)
		}
	})

	r.Post("/within", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		var req WithinRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			panic(err)
		}

		version := semver.MustParse(req.Version)
		spec, err := semver.NewConstraint(req.Spec)
		if err != nil {
			panic(err)
		}

		_, err = w.Write([]byte(strconv.FormatBool(spec.Check(version))))
		if err != nil {
			panic(err)
		}
	})

	r.Get("/versions", func(w http.ResponseWriter, r *http.Request) {
		fixtures, err := bindata.Asset("../../bindata/data/versions.json")
		if err != nil {
			panic(err)
		}

		_, err = w.Write(fixtures)
		if err != nil {
			panic(err)
		}
	})

	// r.Get("/manifests", func(w http.ResponseWriter, r *http.Request) {})

	p := strconv.Itoa(*port)
	fmt.Printf("Listening to :%s\n", p)
	http.ListenAndServe(":"+p, r)
}
