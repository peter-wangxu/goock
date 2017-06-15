# Goock

[![Go Report Card](https://goreportcard.com/badge/github.com/peter-wangxu/goock)](https://goreportcard.com/report/github.com/peter-wangxu/goock)
[![CircleCI](https://img.shields.io/circleci/project/github/peter-wangxu/goock/master.svg?style=flat-square)](https://circleci.com/gh/peter-wangxu/goock)
[![Codecov](https://img.shields.io/codecov/c/github/peter-wangxu/goock/master.svg?style=flat-square)](https://codecov.io/gh/peter-wangxu/goock)

Goock is a `Golang` library/client for discovering and managing block device. it dramatically eases the
efforts needed when connecting/disconnecting to storage backend.

## Table of Content

* [Overview](#overview)
* [Features](#features)
* [Installation](#installation)
* [Requirements](#requirements)
* [Usage](#usage)
    * [As a library](#as-a-library)
        * [Connect to a storage device](#connect-to-a-storage-device)
        * [Disconnect a device from storage system](#disconnect-a-device-from-storage-system)
    * [As a client](#as-a-client)
        * [Show the goock version](#show-the-goock-version)
        * [Connect to a LUN on specific target](#connect-to-a-lun-on-specific-target)
        * [Connect and rescan all LUNs from a target](#connect-and-rescan-all-luns-from-a-target)
        * [Disconnect a device from remote system](#disconnect-a-lun-from-storage-system)
        * [Extend a connected device](#extend-a-connected-device)
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

Note: if you want build binary from source, please firstly [setup Go environment](https://golang.org/doc/).

- Download the source and it's dependencies from github

```
go get -d -v -t github.com/peter-wangxu/goock
```
- Build the binary

```bash
go build
```
a binary file named `goock` will be in place for use.

- Install tools

This step installs couple of tools that goock relies on. 

On Debian/Ubuntu
```bash
sudo apt-get install open-iscsi multipath-tools sysfsutils sg3-utils
```
On RHEL/CentOS

```bash
yum install iscsi-initiator-utils device-mapper-multipath sysfsutils sg3_utils
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

#### Disconnect a device from storage system

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

#### Extend a already connected device

Sometimes, the device can be extend on the storage system, while the device size
is not awared from the host side, in this case, a host side rescan is needed.

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

deviceInfo, _ : = iscsi.ExtendVolume(conn)
```

### As a client tool

Goock is also client tool, which can be used from shell. When the host is connecting with
backend storage(such as enterprise array, DellEMC Unity, HPE 3Par etc.), goock can be helpful
to connect with the storage device and inspect its information.


NOTE: make sure you have permission to operate the block device.


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

#### Disconnect a LUN from storage system
 
 ```bash
 goock disconnect <Target IP> <LUN ID>
 ```

#### Extend a connected device

```bash
goock extend <Target IP> <LUN ID>
```
or

```bash
goock extend /dev/sdx
```

#### Get help

```bash
goock help connect
```
## Testing

### Unit test
```bash
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
