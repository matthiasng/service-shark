# Service Shark

![GitHub](https://img.shields.io/github/license/matthiasng/service-shark)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/matthiasng/service-shark?sort=semver)
![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/matthiasng/service-shark/build/master)
[![codecov](https://codecov.io/gh/matthiasng/service-shark/branch/master/graph/badge.svg)](https://codecov.io/gh/matthiasng/service-shark)
[![Go Report Card](https://goreportcard.com/badge/github.com/matthiasng/service-wrapper)](https://goreportcard.com/report/github.com/matthiasng/service-wrapper)

Service Shark can be used to to host any executable as an Windows service.

Service Shark is:
- easy to use
- lightweight (~2 MB)
- has zero runtime dependencies (no .NET Framework, Java, ...)
- writte in [golang](https://golang.org/)

Service Shark is not:
- a service manager. Their are already easy ways to manage Windows services ([powershell](https://docs.microsoft.com/de-de/powershell/scripting/samples/managing-services?view=powershell-6), [cmd](https://docs.microsoft.com/de-de/windows-server/administration/windows-commands/sc-create), [NSSM](https://nssm.cc/))

## Installation

### Pre-compiled binary

#### Download from github
```
https://github.com/matthiasng/service-shark/releases/latest)
```

#### Scoop
```
#todo setup scoop bucket + goreleaser
```

### Compiling from source
```
git clone https://github.com/matthiasng/service-shark.git
cd service-shark
go build -o service-shark.exe main.go
```

## Usage

#todo
