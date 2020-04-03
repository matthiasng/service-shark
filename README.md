# Service Shark

![GitHub](https://img.shields.io/github/license/matthiasng/service-shark)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/matthiasng/service-shark?sort=semver)
![build](https://github.com/matthiasng/service-shark/workflows/build/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/matthiasng/service-shark/branch/master/graph/badge.svg)](https://codecov.io/gh/matthiasng/service-shark)
[![Go Report Card](https://goreportcard.com/badge/github.com/matthiasng/service-shark)](https://goreportcard.com/report/github.com/matthiasng/service-shark)

Service Shark can be used to to host any executable as a Windows service.

Service Shark is:
- easy to use
- lightweight (~2 MB)
- has zero runtime dependencies (no .NET Framework, Java, ...)
- [12factor/config](https://12factor.net/config) support
- written in [golang](https://golang.org/)

Service Shark is not:
- a service manager. There are already ways to manage Windows services ([powershell](https://docs.microsoft.com/de-de/powershell/scripting/samples/managing-services?view=powershell-6), [cmd](https://docs.microsoft.com/de-de/windows-server/administration/windows-commands/sc-create), [NSSM](https://nssm.cc/))

## Installation

### Pre-compiled binary

#### Download from github
```
https://github.com/matthiasng/service-shark/releases/latest
```

### Compiling from source
```
git clone https://github.com/matthiasng/service-shark.git
cd service-shark
go build -o service-shark.exe main.go
```

## Usage

```
  -name string
        Service name [required]
  -workdir string
        Working directory (default "./")
  -logdir string
        Log directory.
        File name: {name}_YYYY-MM-DD_HH-MM-SS (default "./log")
  -cmd string
        Command [required]
  -version
        Print version and exit
  -- (terminator)
        Pass all arguments after the terminator "--" to the command.
        Bind argument to environment variable with "env:{VAR_NAME}".
```

Example
```
service-shark.exe -name MyService -workdir C:/MyService -cmd java -- -jar MyProg.jar -Xmx1G -myArg "env:MY_ENV_VALUE"
```
Service Shark will run ``java`` with ``-jar MyProg.jar -Xmx1G -myArg "123"`` from ``C:/MyService`` (``MY_ENV_VALUE`` is ``123``).

See [example/test-example-service.ps1](./example/test-example-service.ps1) for a complete example.
