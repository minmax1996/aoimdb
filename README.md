# aoimdb
[![Go](https://github.com/minmax1996/aoimdb/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/minmax1996/aoimdb/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/minmax1996/aoimdb/branch/main/graph/badge.svg)](https://codecov.io/gh/minmax1996/aoimdb)

Another One In Memory Database

## Server
[Server application Readme](https://github.com/minmax1996/aoimdb/blob/main/cmd/aoimd/README.md)
## Client
[simple cli tool](https://github.com/minmax1996/aoimdb/blob/main/cmd/aoimd-cli/README.md)

## Try It

```bash
go get -u github.com/minmax1996/aoimdb/cmd/{aoimd,aoimd-cli}
```

Start server app

```bash
./aoimd
```

Connect with cli tool and try it

```bash
./aoimd-cli -u admin -p pass
```
