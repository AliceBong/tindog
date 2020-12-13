package main

import (
	"github.com/barnslig/torture/lib/elastic"
	"log"
	"time"
)

type ModelFileServerEntry struct {
	Url  string
	Path string
}

type ModelFileEntry stru