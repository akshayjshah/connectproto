connectproto
============

[![Build](https://github.com/akshayjshah/connectproto/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/akshayjshah/connectproto/actions/workflows/ci.yaml)
[![Report Card](https://goreportcard.com/badge/go.akshayshah.org/connectproto)](https://goreportcard.com/report/go.akshayshah.org/connectproto)
[![GoDoc](https://pkg.go.dev/badge/go.akshayshah.org/connectproto.svg)](https://pkg.go.dev/go.akshayshah.org/connectproto)


`connectproto` allows users of [`connect-go`][connect-go] to customize the
default JSON codecs.

## Installation

```
go get go.akshayshah.org/connectproto
```

## Usage

Use this package's options just like the ones built into `connect-go`:

```go
opt := connectproto.WithJSON(
  protojson.MarshalOptions{UseProtoNames: true},
  protojson.UnmarshalOptions{DiscardUnknown: true},
)

// The pingv1connect package is generated from your Protocol Buffer schemas
// by protoc-gen-connect-go. You can use connectproto options with 
// both handlers and clients.
route, handler := pingv1connect.NewPingServiceHandler(
  &pingv1connect.UnimplementedPingServiceHandler{},
  opt,
)
client := pingv1connect.NewPingServiceClient(
  http.DefaultClient,
  "https://localhost:8080",
  opt,
)
```

## Status: Unstable

This module is unstable, with a stable release expected before the end of 2023.
It supports:

* The [two most recent major releases][go-support-policy] of Go.
* [APIv2] of Protocol Buffers in Go (`google.golang.org/protobuf`).

Within those parameters, `connectproto` follows semantic versioning. 

## Legal

Offered under the [MIT][license].

[APIv2]: https://blog.golang.org/protobuf-apiv2
[go-support-policy]: https://golang.org/doc/devel/release#policy
[license]: https://github.com/akshayjshah/connectproto/blob/main/LICENSE
[connect-go]: github.com/bufbuild/connect-go
