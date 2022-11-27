# Aoimd

Another One In Memory Database

this tool for setting up database in memory, settle up with endpoints and listen for clients

by default it listen `:1593` port for tcp-connects
It can be accessed with [simple cli tool](https://github.com/minmax1996/aoimdb/blob/main/cmd/aoimd-cli/README.md)

data encoding between cli client by default `msgpack`

## DataTypes

1) SET ( like redis one)
simple KV implemented by `map[string]interface{}`;
interact with `set` and `get` commands
2) ...

## DataBackups

Every 30 sec server will write data to file; On start it will try to restore data from backup

## Future Features

* Web Interface to access data