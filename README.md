[![Build Status](https://travis-ci.org/sebnow/httptracing.svg?branch=master)](https://travis-ci.org/sebnow/httptracing)
[![GoDoc](https://godoc.org/github.com/sebnow/httptracing?status.svg)](http://godoc.org/github.com/sebnow/httptracing)
[![Coverage Status](https://coveralls.io/repos/github/sebnow/httptracing/badge.svg?branch=master)](https://coveralls.io/github/sebnow/httptracing?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sebnow/httptracing)](https://goreportcard.com/report/github.com/sebnow/httptracing)

Description
===========

The `httptracing` package forms a bridge between [httpclient] and the
opentracing-contrib [nethttp] package. It provides a HTTP client with
opentracing support, implementing the `httpclient.Client` interface.

[httpclient]: https://github.com/sebnow/httpclient
[nethttp]: https://github.com/opentracing-contrib/go-stdlib


Usage
=====

```go
tracer := opentracing.GlobalTracer()
client := httptracing.Trace(tracer, http.DefaultClient)

resp, err := client.Get("http://test.com")
//...
```
