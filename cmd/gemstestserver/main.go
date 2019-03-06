package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

//go:generate go-bindata -o ../../bindata/bindata.go -pkg bindata ../../bindata

func main() {
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
		fixtures, err := bindata.Asset("../../bindata/versions.json")
		if err != nil {
			panic(err)
		}

		_, err = w.Write(fixtures)
		if err != nil {
			panic(err)
		}
	})

	// r.Get("/manifests", func(w http.ResponseWriter, r *http.Request) {})

	http.ListenAndServe(":3333", r)
}
