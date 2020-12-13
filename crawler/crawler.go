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

	// Initially load config
	crawlers.Load(configPath)

	// Initialize config reloading using syscall
	reloadChan := make(chan os.Signal)
	signal.Notify(reloadChan, syscall.SIGUSR1)
	go func() {
		for {
			<-reloadChan

			log.Println("reload!")
			crawlers.Load(configPath)
		}
	}()

	return
}

// (Re-)load crawler configs and create crawler instances. You can call this
// function again multiple times to reload the configs
func (crawlers *Crawlers) Load(configPath string) (err error) {
	rawConfig, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return
	}

	nextConfig := CrawlersConfig{}
	err = json.Unmarshal(rawConfig, &nextConfig)
	if err != nil {
		return
	}

	var nextCrawlers []*CrawlerEntry
	for _, entrypoint := range nextConfig.Entrypoints {
		// Parse config while providing default values
		entryConfig := CrawlerConfig{
			TurnDelay: 10 * time.Second,
		}
		err = json.Unmarshal(*entrypoint, &entryConfig)
		if err != nil {
			return
		}

		nextCrawlers = append(nextCrawlers, &CrawlerEntry{
			Config:    entryConfig,
			RawConfig: entrypoint,
		})
	}

	// 1. check for removed crawlers and terminate them
OUTER:
	for _, oldCrawler := range crawlers.Crawlers {
		for _, newCrawler := range nextCrawlers {
			if bytes.Compare([]byte(*oldCrawler.RawConfig), []byte(*newCrawler.RawConfig)) == 0 {
				// if crawler config has not changed, keep the crawler instance and continue
				newCrawler.Crawler = oldCrawler.Crawler
				continue OUTER
			}
		}

		// otherwise terminate the old crawler
		log.Printf("server %s terminated\n", oldCr