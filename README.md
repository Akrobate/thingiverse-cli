# thingiverse-cli

thingiverse-cli command line tool. Allows you to upload directly your updates to thingiverse or can be used as a part of the CI for publishing automaticly updates


## Commands

### Configuration

The first command to enter for setting the client id and client secret. This command will create the thingiverse-cli configuration file on your system. Should be called once to configure the account.

```bash
thingiverse-cli config
```

### Authentication

Launch the authentication process. Just follow the steps, you will be asked to visit an url, and your browser will print you the access token. You'll have to copy paste the token in command line for saving it

```bash
thingiverse-cli auth
```

### init

The init command create an empty `thingiverse.yml` file in your project, once created you'll need to fill the fields with your data

```bash
thingiverse-cli init
```

### View categories referential

Show the possible values for categories

```bash
thingiverse-cli remote categories
```

### search tags

Search the existing tags values and counts

```bash
thingiverse-cli remote tags search_string
```

### View licenses referential

```bash
thingiverse-cli remote licenses
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

## Example thingiverse.yml

```yaml
iid: 7387159
name: TP4056 LiPo Charger Panel / Surface Mount Bracket
category: 66
license: cc
is_wip: false
tags:
  - Tp4056
  - tp4056_case
  - tp4056_enclosure
  - TP4056_Charger_Module
  - TP4056_charging_board
image_files:
  - local_path: ./assets/photos/printed_preview_2.jpg
  - local_path: ./assets/photos/printed_preview_1.jpg
  - local_path: ./opm_png_files/example.png
  - local_path: ./opm_png_files/example-exploded.png
model_files:
  - local_path: ./opm_stl_files/tp4056FacadeHolderPiece.stl
  - local_path: ./opm_stl_files/tp4056FacadeRoundedHolderPiece.stl
  - local_path: ./opm_stl_files/tp4056FacadeSquareMiniHolderPiece.stl
  - local_path: ./opm_stl_files/tp4056FixationPiece.stl
instructions: ""
description: |
  This is a modular, 2-piece mounting bracket designed to securely hold a standard **TP4056 LiPo battery charging module** in place. Whether you are embedding a charging port into a custom project enclosure (panel mount) or securing it onto a flat surface, this mount keeps your board aligned and protected.

  The design relies on a clamp system held together by **two M3 screws and M3 nuts** to provide a tight, strain-relieved hold on the PCB without needing glue.

  ### ⚙️ Features & Hardware Requirements

  * **Mounting Style:** Panel mount (cutout) or surface mount.
  * **Compatibility:** Standard TP4056 USB-C / Micro-USB charging boards.
  * **Hardware Required:** 
    * 2x M3 screws (M3x8mm to M3x12mm recommended)
    * 2x M3 nuts

  ### 🔓 Source Code & Customization

  The original **OpenSCAD source code** for this model is open-source and hosted on GitLab. If you want to modify dimensions or customize the design:

  * **Git repository:** [tp4056-holder](https://gitlab.com/openscad-modules/tp4056-holder)

  Feel free to fork, customize, or contribute back!

  ### 📦 OpenSCAD Dependencies & Package Management

  This model uses external modules managed with **[opm (OpenSCAD Package Manager)](https://openscad-modules.gitlab.io/openscad-package-manager-documentation/)**. 

  To easily build and customize this model with all required libraries:

  1. **Install `opm`**: Follow the [OpenSCAD Package Manager Quick Start Guide](https://openscad-modules.gitlab.io/openscad-package-manager-documentation/quick-start/).
  2. **Download Dependencies**: Open your terminal in the project directory and run:
    ``bash
    opm install
    ``

  ### 🧩 Use as an OpenSCAD Module

  This project can also be imported directly as a reusable module in your own OpenSCAD designs via **opm** (the **OpenSCAD Package Manager**).

  To install this package in your project, run:

  ``bash
  opm install https://gitlab.com/openscad-modules/tp4056-holder.git#0.0.1
  ``
```

## Todo

- [X] Init empty thingiverse.yml File command line
- [X] Search tags command line
- [ ] Reorder galleries images
- [ ] Reorder Files