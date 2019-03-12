# gemstest

`gemstest` provides files to support laptop-based implementation for the Gems
interview challenge.

## Installation

Download from the GitHub Releases list. Use https://transfer.sh to download to
a candidate's computer.

## Running the server

To run: `gemstestserver`

For usage: `gemstestserver --help`

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

## Generating data files

Download a RubyGems database dump from https://rubygems.org/pages/data, and load
it into the provided Docker container. Then run the Docker container with its
PostgreSQL port exposed and run `gengems`.
