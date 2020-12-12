# torture
Multi FTP crawler and file search. Written for CCC events.

torture exists of two seperate components: The crawler and the frontend. Both are interconnected via the Elasticsearch backend.

## Setup
It is always a good idea to setup a [GOPATH](https://golang.org/doc/code.html#GOPATH).

Dependencies are managed using [dep](https://github.com/golang/dep).

### General
1. Install and setup [Elasticsearch](https: