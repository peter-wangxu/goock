# Goock

[![CircleCI](https://img.shields.io/circleci/project/github/peter-wangxu/goock/master.svg?style=plastic)]()

[![Codecov](https://img.shields.io/codecov/c/github/peter-wangxu/goock/master.svg?style=plastic)]()

-----

Goock is a go library for discovering and managing block device.

## Introduction

Goock aims at providing a easy-to-use library for block device discovering and
management. This project is inspired by OpenStack project
[os-brick](https://github.com/openstack/os-brick).

### Features

* Discovery of devices for iSCSI transport protocol.
* Discovery of devices for FibreChanel transport protocol.
* Removal of devices from a host
* Multipath support

## Installation

```
go get github.com/peter-wangxu/goock
```

## Requirements

* Linux
* Go 1.7 or later

## Testing


```
cd goock
go test -v ./...
```

## Contributions

Simply fork this repo and send PR for your code change(also tests to cover your
change), remember to give a title and description of your PR.

## License
[Apache License Version 2.0](LICENSE)

## FAQ
