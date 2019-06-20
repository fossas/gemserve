# gemstest

`gemstestserver` provides files to support laptop-based implementation for the Gems interview challenge.

## Installation

Download from the GitHub Releases list. Use https://transfer.sh to download to a candidate's computer, or host on your own machine using https://ngrok.io.

## Running the server

To run: `gemstestserver`

For usage: `gemstestserver --help`

## Usage

This HTTP server listens on three paths:

### `GET /versions`

Returns a JSON list containing a list of (package name, package version) pairs.

Example output:

```
[
  ["package", "1.0.0"],
  ["package", "1.0.1"],
  ["another-package", "2.3.4"]
]
```

### `GET /manifests`

Returns a JSON object containing package manifests.

Example output:

```
{
  "package": {
    "a-direct-dependency": "~> 1.0.0",
    "yet-another-package": "^2.3.4"
  },
  "another-package": {
    "yet-another-package": ">= 2, < 4.0.0"
  }
  "yet-another-package": {}
}
```

### `POST /compare`

Expects a JSON request body with two strings "a" and "b".
Returns -1 if a < b, 0 if a == b, and 1 if a > b.

Example input:

```
{
  "a": "1.0.0",
  "b": "1.0.1"
}
```

Example output:

```
-1
```

### `POST /within`

Expects a JSON request body with two strings "version" and "spec". Returns
true if the version is within the spec, and false otherwise.

Example input:

```
{
  "version": "1.0.2",
  "spec": "^1.0.0"
}
```

Example output:

```
true
```

## Development

### Generating data files

Download a RubyGems database dump from https://rubygems.org/pages/data, and load
it into the provided Docker container. Then run the Docker container with its
PostgreSQL port exposed and run `gengems`.
