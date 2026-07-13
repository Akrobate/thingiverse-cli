# thingiverse-cli
thingiverse-cli command line tool for upload to thingiverse

## Prototyping thingiverse.json

```json
{
  "thingId": "1234567",
  "name": "tp4056 Holder - opm module",
  "license": "cc-by-nc-sa",
  "category": "3d-printing",
  "files": [
    {
      "localPath": "./stl_files/wing_left.stl",
      "name": "wing_left_v2.stl"
    },
    {
      "localPath": "./stl_files/wing_right.stl",
      "name": "wing_right_v2.stl"
    }
  ]
}
```

## Static page for retrieve token

https://akrobate.github.io/thingiverse-cli/token.html


## Development

### Requirements

- Go 1.23 or newer

### Build

```bash
go build -o thingiverse-cli
```

### Build and install locally

```bash
go build -o thingiverse-cli && sudo cp thingiverse-cli /usr/local/bin/
```