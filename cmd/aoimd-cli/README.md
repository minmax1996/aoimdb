# Aoimd

Another One In Memory Database CLI

## Description

Simple client for interacting with [server application](https://github.com/minmax1996/aoimdb/blob/main/cmd/aoimd/README.md)

## Usage

`aoimd-cli`

### Flags

By default it will try to connect to `localhost:1593` server instance

* `-h` string flag for host
* `-u` string flag for user
* `-p` string flag for pass

Examples:

`aoimd-cli -h localhost:1593`

`aoimd-cli -u admin -p pass`

`aoimd-cli -u admin`

### Commands

After connection `aoimd-cli` sends `help` command by default
```
Usage:
<command> <args>

Commands:
        help     shows this message
        auth     auth database command (Usage: auth user pass)
        select   select database command (Usage: select <database_name>)
        get      get command (Usage: get [<databasename>.]<key>)
        set      set database command (Usage: set [<databasename>.]<key> <value>)
        keys     keys database command (Usage: keys [keysregexp])
        exit     exits from server
```

command `help` with arg will print arg help

```
help set
set      set database command (Usage: set [<databasename>.]<key> <value>
```
