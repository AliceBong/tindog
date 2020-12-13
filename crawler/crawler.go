package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"
)

// Data structures shared by all protocol-specific crawlers

type FileInfo struct {
	URL      *url.URL
	Size     int64
	MimeType string
	ModTime  time.Time
}

type WalkFunction func(path string, info FileInfo)

type Crawler interface {
	// Start walking recursively and call fn on every file
	Walk(fn WalkFunction) error

	// Tear down all open connections
	Close()
}

// Data structures only used to control instances of protocol specific crawlers

type CrawlerConfig struct {
	Entry     string        `json:"entry"`
	TurnDelay time.Duration `json:"turnDelay"`
}

type CrawlersConfig struct {
	Entrypoints []*json.RawMessage `json:"entrypoints"`
}

type CrawlerEntry struct {
	Config    CrawlerConfig
	RawConfig *json.RawMessage
	Crawler   *Crawler
	Terminate chan bool
}

type Crawlers struct {
	Config    CrawlersConfig
	Crawlers  []*CrawlerEntry
	WaitGroup sync.WaitGroup
	Model     *Model
}

func CreateCrawlers(configPath string, model *Model) (crawlers *Crawlers, err error) {
	crawlers = &Crawlers{
		Model: model,
	}
