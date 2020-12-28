![goreleaser](https://github.com/eqinox76/gotestor/workflows/goreleaser/badge.svg)

# gotestor
Simple tool to continously run golang unit tests.

# install
Either build and install it directly via
```
go get github.com/eqinox76/gotestor
```
and make sure the $GOROOT/bin or $GOTOOLDIR directory is in your path.

Or use the prebuild [executables](https://github.com/eqinox76/gotestor/releases/).

# usage
Navigate to your code directory and start gotestor it will immediatly run all lcoal unit tests
and start listening for file changes.
```
gotestor
2020/12/24 15:04:30 running  /snap/bin/go test
?   	github.com/eqinox76/gotestor	[no test files]
2020/12/24 15:04:53 running  /snap/bin/go test
?   	github.com/eqinox76/gotestor	[no test files]
```

You can pass the usual go test parameters directly
```
gotestor -count=1 -failfast
```

# intention
This tool supports my personal workflow of changing code and running all unit tests immediatly.
It is a very ligth wrapper around the `go test` tool and only adds the automatic run if any source files changes.
All parameters are forwarded.

# alternatives
There are several ways to solve this issue with other tools. Starting from simple bash loops up to sophisticated npm packages.

# thanks
* [fsnotify](https://github.com/fsnotify/fsnotify)
