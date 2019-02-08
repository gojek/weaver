# Weaver

<a href="https://travis-ci.org/gojektech/heimdall"><img src="https://travis-ci.org/gojektech/heimdall.svg?branch=master" alt="Build Status"></img></a> [![Go Report Card](https://goreportcard.com/badge/github.com/gojekfarm/weaver)](https://goreportcard.com/report/github.com/gojekfarm/weaver)

* [Description](#description)
* [Installation](#installation)
* [Usage](#usage)
* [Documentation](#documentation)
* [FAQ](#faq)
* [License](#license)

## Description
Weaver is a Layer-7 Load Balancer with Dynamic Sharding Strategies. 
It is a simple HTTP reverse proxy.

## Features:

- Sharding request based on headers/path/body fields
- Emits Metrics on requests per route per backend
- Dynamic configuring of different routes (No restarts!)
- Is Fast
- Supports multiple algorithms for sharding requests (consistent hashing, modulo, s2 etc)
- Packaged as a single self contained binary
- Logs on failures (Observability)

## Installation

### Build from source

- Clone the repo:
```
git clone git@github.com:gojektech/weaver.git
```

- Build to create weaver binary
```
make build
```

### Binaries for various architectures

Download the binary for a release from: [here](https://github.com/gojekfarm/weaver/releases)

## License

```
Copyright 2018, GO-JEK Tech (http://gojek.tech)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
