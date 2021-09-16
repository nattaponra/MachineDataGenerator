# MachineDataGenerator

![Tag](https://img.shields.io/github/v/tag/ifraiot/MachineDataGenerator)
![Last commit date](https://img.shields.io/github/last-commit/ifraiot/MachineDataGenerator)
![Go version](https://img.shields.io/github/go-mod/go-version/ifraiot/MachineDataGenerator)
![Repo size](https://img.shields.io/github/repo-size/ifraiot/MachineDataGenerator)
![Contributor](https://img.shields.io/github/contributors/ifraiot/MachineDataGenerator)

Project URL: [https://github.com/ifraiot/MachineDataGenerator](https://github.com/ifraiot/MachineDataGenerator)

By Thitiwut Chutipongwanit

MachineDataGenerator helps to publish MQTT message perpetually from multiple publishers.

## Prerequisites

Before installing `MachineDataGenerator` you need:

- Git
- Go 1.16+

## Quick start

Set environment variables in .env file.

```shell
cat .env
```

```text
MDG_BROKER_URL=<MQTT BROKER URL>
```

Set publisher configurations.

```shell
cat data/publishers.json
```

```json
[
  {
    "username": "<PUBLISHER USERNAME (deviceId)>",
    "password": "<PUBLISHER PASSWORD (primaryKey)>",
    "jsonOpPath": "<JSON OUTPUT MESSAGES FILEPATH>",
    "jsonRopPath": "<JSON REJECT OUTPUT MESSAGES FILEPATH>",
    "jsonStPath": "<JSON STATUS MESSAGES FILEPATH>",
    "topic": "<MQTT TOPIC TO PUBLISH TO>"
  }
]
```

Run MachineDataGenerator

```bash
go run main.go
```

## Usage

```shell
$ machine-data-generator --help
Usage: machine-data-generator 
```

## Configuration

### publishers.json

`publishers.json` must be located at
`<PROJECT ROOT DIRECTORY>/data/publishers.json` in order for
MachineDataGenerator to be able to work properly.

Format

```json
[
  {
    "username": "<PUBLISHER USERNAME (deviceId)>",
    "password": "<PUBLISHER PASSWORD (primaryKey)>",
    "jsonOpPath": "<JSON OUTPUT MESSAGES FILEPATH>",
    "jsonRopPath": "<JSON REJECT OUTPUT MESSAGES FILEPATH>",
    "jsonStPath": "<JSON STATUS MESSAGES FILEPATH>",
    "topic": "<MQTT TOPIC TO PUBLISH TO>"
  }
]
```

## Environment variables

MachineDataGenerator use shell environment variable(s).

Variables can be set in `.env` file at root project directory.
If `.env` file is found, MachineDataGenerator load it to shell environment.

`MDG_BROKER_URL`

Source MQTT broker URL.

The format should be `scheme://host:port`Where
`scheme` is one of "tcp", "ssl", or "ws",
`host` is the ip-address (or hostname) and
`port` is the port on which the broker is accepting connections.

## Support

Thitiwut Chutipongwanit - thitiwut@ifrasoft.com
