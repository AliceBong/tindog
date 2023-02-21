package elastic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"log"
)

type Hit struct {
	Id     string           `json:"_id"`
	Score  float32          `json:"_score"`
	Source *json.RawMessage `json:"_source"`
}

type Hits struct {
	Total int   `json:"total"`
	Hits  []Hit `json:"hits"`
}

type Result struct {
	Hits         Hits             `json:"hits"`
	Aggregations *j