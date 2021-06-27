# Open Keyless
[![Build Status](https://cloud.drone.io/api/badges/betterengineering/open-keyless/status.svg)](https://cloud.drone.io/betterengineering/open-keyless)
[![GoDoc](https://godoc.org/github.com/betterengineering/open-keyless?status.svg)](https://godoc.org/github.com/betterengineering/open-keyless)

Open Keyless is a keyless entry system for contactless key cards. This project includes all of the necessary software,
printed circuit board designs, 3D models, and documentation to build an RFID based badge reader and controller.

Note: this project is under development and has not yet been completed. Once a functioning version of this
project exists, this message will be removed and the release will be added as a github release.

## Installation
To install the latest release, run the following from your Raspberry Pi:
```
curl -LO https://github.com/betterengineering/open-keyless/releases/latest/download/install.sh | sudo bash
```

## Documentation
The documentation for Open Keyless is kept in the repo! Checkout the [Overview](docs/overview.md) for a starting point.

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
