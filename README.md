# kickstart.go
Minimalistic http server template in go that is:
- Small (less than 300 lines of code)
- Single file 
- Only standard library dependencies

**Not** a framework, but a starting point for building HTTP services in Go.  

Inspired by [Mat Ryer](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/#designing-for-testability) & [earthboundkid](https://blog.carlana.net/post/2023/golang-git-hash-how-to/) and even [kickstart.nvim](https://github.com/nvim-lua/kickstart.nvim)

## Features
- Graceful shutdown: Handles `SIGINT` and `SIGTERM` signals to shutdown gracefully.
- Health endpoint: Returns the server's health status including version and revision.
- OpenAPI endpoint: Serves an OpenAPI specification.
- Debug information: Provides various debug metrics including pprof and expvars.
- Access logging: Logs request details including latency, method, path, status, and bytes written.
- Panic recovery: Catch and log panics in HTTP handlers gracefully.

## Getting started
- Use this template to create a new repository
- Or fork the repository and make changes to suit your needs.

### Requirements
Go 1.22 or later

### Suggested Dependencies
- [golangci-lint](https://golangci-lint.run/) 

### Build and run the server
```sh
$ make run 
```
- this will build the server and run it on port 8080
- Checkout Makefile for more 

## Endpoints
- GET /health: Returns the health of the service, including version, revision, and modification status.
- GET /openapi.yaml: Returns the OpenAPI specification of the service.
- GET /debug/pprof: Returns the pprof debug information.
- GET /debug/vars: Returns the expvars debug information.

## OpenAPI
- The OpenAPI definition file is embedded in the binary using Go's embed package and serves at the /openapi.yaml endpoint. Modify the api/openapi.yaml file to change the OpenAPI specifications.
