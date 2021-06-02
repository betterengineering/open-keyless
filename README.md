# Open Keyless
[![Build Status](https://cloud.drone.io/api/badges/betterengineering/open-keyless/status.svg)](https://cloud.drone.io/betterengineering/open-keyless)
[![GoDoc](https://godoc.org/github.com/betterengineering/open-keyless?status.svg)](https://godoc.org/github.com/betterengineering/open-keyless)

Open Keyless is a keyless entry system for contactless key cards. This project includes all of the necessary software,
printed circuit board designs, 3D models, and documentation to build an RFID based badge reader and controller.

Note: this project is under active development and has not yet been completed. Once a functioning version of this
project exists, this message will be removed and the release will be added as a github release.

## Documentation
The documentation for Open Keyless is kept in the repo! Checkout the [Overview](docs/overview.md) for a starting point.

## Versioning
This projects contains 3D models, printed circuit boards, and software that is released as a unit to guarantee
compatibility. The software release follows [semantic versioning](https://semver.org/). The hardware release also
follows semantic versioning with the exception of there not being a patch version. Because the project is verisoned as
a unit, changes in either the software or the hardware will cause version increments in both.

## Development
This project uses Go 1.11 modules! If you have Go 1.11 or later installed, no other project setup is necessary.

### Testing
To run tests locally, run the following:
```bash
make test
```

To include integration tests, run the following:
```bash
make test-full
```

### Mocks
This project uses Golang's [mockgen](https://github.com/golang/mock) tool to generate mocks for interfaces. To regenerate mocks, run the following:
```bash
make mocks
```

### Libnfc
To develop locally, you will need libnfc installed. I would like to remove this as a dependency in the future, but alas,
it is a dependency for now.

For MacOS, run the following:
```bash
brew install libnfc
```

For Debian based Linux distros, run the following:
```bash
apt-get install -y libnfc-dev
```
