# lookslike

[![Build Status](https://travis-ci.org/elastic/go-lookslike.svg?branch=master)](https://travis-ci.org/elastic/go-lookslike)

This library is here to help you with all your data validation needs. It's ideally suited to writing tests against JSON-like structures, but can do much much more. We use it in [elastic/beats](https://github.com/elastic/beats) for our own tests. 

## Quick Links

* [GoDoc](https://godoc.org/github.com/elastic/go-lookslike) for this library.
* [Runnable Examples](https://github.com/elastic/go-lookslike/blob/master/lookslike/doc_test.go).

## Real World Usage Examples

lookslike was created to improve the testing of various structures in [elastic/beats](https://github.com/elastic/beats). Searching the tests for `lookslike` will show real world usage.

## Call for More `isdef`s!

We currently [define](https://godoc.org/github.com/elastic/go-lookslike/lookslike/isdef) only the isdefs
we've actually used in the field. If you'd like to add your own, please open a PR (with tests!).
