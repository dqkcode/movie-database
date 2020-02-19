package movie

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/dqkcode/movie-database/internal/app/types"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

type (
	Elasticsearch struct {
		client *elasticsearch.Client
	}
)

func NewElasticsearch(es *elasticsearch.Client) *Elasticsearch {
	return &Elasticsearch{
		client: es,
	}
}

var (
	r map[string]interface{}
)

func (s *Elasticsearch) InsertMovies(m *types.MovieInfo) error {
	id := m.ID
	m.ID = ""
	x, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:      "movies",
		DocumentID: id,
		Body:       strings.NewReader(string(x)),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), m.ID)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}
	return nil
}
