# Goock

[![CircleCI](https://img.shields.io/circleci/project/github/peter-wangxu/goock/master.svg?style=plastic)]()
[![Codecov](https://img.shields.io/codecov/c/github/peter-wangxu/goock/master.svg?style=plastic)]()

Goock is a GO library/client for discovering and managing block device. it dramatically eases the
effort needed when connecting to storage backend.

## Table of Content

* [Overview](#overview)
* [Features](#features)
* [Installation](#installation)
* [Requirements](#requirements)
* [Usage](#usage)
    * [As a library](#as-a-library)
        * [Connect to a storage device](#connect-to-a-storage-device)
        * [Disconnect from a storage device](#disconnect-from-a-storage-device)
    * [As a client](#as-a-client)
        * [Show the goock version](#show-the-goock-version)
        * [Connect to a LUN on specific target](#connect-to-a-lun-on-specific-target)
        * [Connect and rescan all LUNs from a target](#connect-and-rescan-all-luns-from-a-target)
        * [Get help](#get-help)
* [Testing](#testing)
    * [Unit test](#unit-test)
    * [Manual test](#manual-test)
* [Contributions](#contributions)
* [License](#license)
* [FAQ](#faq)


## Overview

Goock aims at providing a easy-to-use library/client for block device discovering and management. user no longer needs
the use of multiple linux tools like `open-iscsi`, `multipath-tools`, `sysfsutils` and `sg3-utils`). Instead,  
goock will leverage the power of many *nux device management tools, and provides simple and straightforward API/CLI for
developers and administrators.


This project is inspired by OpenStack project
[os-brick](https://github.com/openstack/os-brick)
## Features

* Discovery of devices for iSCSI transport protocol.
* Discovery of devices for FibreChanel transport protocol.
* Removal of devices from a host
* Multipath support

## Installation

download the source from github

```
go get -d -v -t github.com/peter-wangxu/goock
```

## Requirements

Goock can be built or developed on both Linux or Windows platform.

* Linux/Windows
* Go 1.7 or later

## Usage

### As a library

Goock is a library for connecting/disconnecting block devices for any Golang based
software. the example usage below:

#### Connect to a storage device

```go
package main

import (
        "github.com/peter-wangxu/goock/connector"
)

iscsi := connector.New()

conn := connector.ConnectionProperty{}
conn.TargetPortals = []string{"192.168.1.30"}
conn.TargetIqns = []string{"iqn.xxxxxxxxxxxxxx"}
conn.TargetPortals = []int{10}

deviceInfo, _ : = iscsi.ConnectVolume(conn)
```

#### Disconnect from a storage device

```go
package main

import (
        "github.com/peter-wangxu/goock/connector"
)

iscsi := connector.New()

conn := connector.ConnectionProperty{}
conn.TargetPortals = []string{"192.168.1.30"}
conn.TargetIqns = []string{"iqn.xxxxxxxxxxxxxx"}
conn.TargetPortals = []int{10}

deviceInfo, _ : = iscsi.DisconnectVolume(conn)
```

### As a client tool
Goock is also client tool, which can be used from shell. When the host is connecting with
backend storage(such as enterprise array, DellEMC Unity, HPE 3Par etc.), goock can be helpful
to connect with the storage device and inspect its information.

Example usage below:

#### Show the goock version

```bash
goock -v
```

#### Connect to a LUN on specific target

```bash
# Connect and rescan storage device whose iSCSI target ip is <target IP>, the desired
# Device LUN ID is [LUN ID]
goock connect <target IP> [LUN ID]

```
#### Connect and rescan all LUNs from a target

```bash

goock connect <target IP>

```

#### Get help

```bash
goock help connect
```
## Testing

### Unit test
```
cd goock
go test -v ./...
```

### Manual test

* first build a binary from source

```bash
cd goock
go build
```
* run `goock` command by specifying some parameters
```bash
goock connect <target IP>
```


## Contributions

Simply fork this repo and send PR for your code change(also tests to cover your
change), remember to give a title and description of your PR.

## License

[Apache License Version 2.0](LICENSE)

## FAQ
