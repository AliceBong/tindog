What
====

This is a robots.txt exclusion protocol implementation for Go language (golang).


Build
=====

To build and run tests run `script/test` in source directory.


Contribute
==========

Warm welcome.

* If desired, add your name in README.rst, section Who.
* Run `script/test && script/clean && echo ok`
* You can ignore linter warnings, but everything else must pass.
* Send your change as pull request or just a regular patch to current maintainer (see section Who).

Thank you.


Usage
=====

As usual, no special installation is required, just

    import "github.com/temoto/robotstxt"

run `go get` and you're ready.

1. Parse
^^^^^^^^

First of all, you need to parse robots.txt data. You can do it with
functions `FromBytes(body []byte) (*RobotsData, error)` or same for `string`::

    robots, err := robotstxt.FromBytes([]byte("User-agent: *\nDisallow:"))
    robots, err := robotstxt.FromString("User-agent: *\nDisallow:")

As of 2012-10-03, `FromBytes` is the most efficient method, everything else
is a wrapper for this core function.

There are few convenient constructors for various purposes:

* `FromResponse(*http.Response) (*RobotsData, error)` to init robots data
from HTTP response. It *does not* call `response.Body.Close()`::

    robots, err := robotstxt.FromResponse(resp)
    resp.Body.Close()
    if err != nil {
        log.Println("Error parsing robots.txt:", err.Error())
    }

* `FromStatusAndBytes(statusCode int, body []byte) (*RobotsData, error)` or
`FromStatusAndString` if you prefer to read bytes (string) your