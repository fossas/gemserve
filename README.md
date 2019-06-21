# gemserve

`gemserve` provides files to support laptop-based implementation for the Gems interview challenge.

## Installation

Download from the GitHub Releases list. Use https://transfer.sh to download to a candidate's computer, or host on your own machine using https://ngrok.io.

## Running the server

To run: `gemserve`

For usage: `gemserve --help`

## Usage

This HTTP server listens on three paths:

### `GET /versions`

Returns a JSON list containing a list of (package name, package version) pairs.

Example output:

```json
[
  ["package", "1.0.0"],
  ["package", "1.0.1"],
  ["another-package", "2.3.4"]
]
```

### `GET /manifests`

Returns a JSON object containing package manifests.

Example output:

```json
{
  "package": {
    "a-direct-dependency": "~> 1.0.0",
    "yet-another-package": "^2.3.4"
  },
  "another-package": {
    "yet-another-package": ">= 2, < 4.0.0"
  },
  "yet-another-package": {}
}
```

### `POST /compare`

Expects a JSON request body with two strings "a" and "b".
Returns -1 if a < b, 0 if a == b, and 1 if a > b.

Example input:

```json
{
  "a": "1.0.0",
  "b": "1.0.1"
}
```

Example output:

```json
-1
```

### `POST /within`

Expects a JSON request body with two strings "version" and "spec". Returns
true if the version is within the spec, and false otherwise.

Example input:

```json
{
  "version": "1.0.2",
  "spec": "^1.0.0"
}
```

Example output:

```json
true
```

## Development

### Generating data files

Download a RubyGems database dump from https://rubygems.org/pages/data, and load
it into the provided Docker container. Then run the Docker container with its
PostgreSQL port exposed and run `gengems`.

```bash
docker build -t gemserve .
docker run -p 5432:5432 -it gemserve

# Within Docker container
/docker-entrypoint.sh postgres &
./load-pg-dump -c gems.tar

# Outside of container, while container is still running
gengems -manifestFile ./bindata/data/manifests.json -versionFile ./bindata/data/versions.json > gengems.log
go install ./...
```