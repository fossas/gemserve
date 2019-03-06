# gemstest

`gemstest` provides files to support laptop-based implementation for the Gems
interview challenge.

## Installation
`GO111MODULE=on go get`

## Running server
`go run cmd/gemstestserver/main.go`
## Sample requests
```
curl -X POST http://localhost:8000/compare -d '{ "A":"1.2.3", "B": "1.2.3" }'
```
```
curl -X POST http://localhost:8000/within -d '{"Version":"1.2.3","Spec": "^1.2.1"}'
```
```
curl http://localhost:8000/versions
```

