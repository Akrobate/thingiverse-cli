# thingiverse-cli

thingiverse-cli command line tool for upload to thingiverse


## Commands

### init

```bash
thingiverse-cli init
```

## Prototyping thingiverse.yml

```yaml
id: 7387159
name: TP4056 LiPo Charger Panel / Surface Mount Bracket
category: 3d-printing
license: cc
is_wip: false
tags: []
files:
  - local_path: ./opm_png_files/example.png
  - local_path: ./opm_png_files/example-exploded.png
instructions: ""
description: |
  This is a modular, 2-piece mounting bracket designed to securely hold a standard **TP4056 LiPo battery charging module** in place. Whether you are embedding a charging port into a custom project enclosure (panel mount) or securing it onto a flat surface, this mount keeps your board aligned and protected.

  The design relies on a clamp system held together by **two M3 screws and M3 nuts** to provide a tight, strain-relieved hold on the PCB without needing glue.
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


## Todo

- [ ] Search tags command line