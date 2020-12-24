# gotestor
Simple tool to continously run golang unit tests.

# usage

# intention
This tool supports my personal workflow of changing code and running all unit tests immediatly.
It is a very ligth wrapper around the `go test` tool and only adds the automatic run if any source files changes.
All parameters are forwarded.

# alternatives
There are several ways to solve this issue with other tools. Starting from simple bash loops up to sophisticated npm packages.

# thanks
* [fsnotify](https://github.com/fsnotify/fsnotify)
