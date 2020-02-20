package movie

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/sirupsen/logrus"

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

func (s *Elasticsearch) InsertMovies(ctx context.Context, m *types.MovieInfo) error {
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
	}
	res, err := req.Do(context.Background(), s.client)
	if err != nil {
		logrus.Errorf("Error getting response: %s", err)
		return ErrInsertMovieToESFailed
	}
	if res.IsError() {
		logrus.Errorf("[%s] Error indexing document ID=%s", res.Status(), id)
		return ErrIndexingDocumentFailed
	}

	// Deserialize the response into a map.
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logrus.Infof("Error parsing the response body: %s", err)
	} else {
		// Print the response status and indexed document version.
		logrus.Infof("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	}

	return nil
}

type (
	Query struct {
		Match struct {
			Name string `json:"name"`
		} `json:"match"`
	}

	SearchQuery struct {
		Query `json:"query"`
	}
	ResponseElasticsearch struct {
		Took    int  `json:"took"`
		TimeOut bool `json:"timed_out"`
		Shards  struct {
			ToTal      int `json:"total"`
			Successful int `json:"successful"`
			Skipped    int `json:"skipped"`
			Failed     int `json:"failed"`
		} `json:"_shards"`
		Hits struct {
			ToTal struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
			MaxScore float64 `json:"max_score"`
			Hits     []struct {
				Index  string          `json:"_index"`
				ID     string          `json:"_id"`
				Score  float64         `json:"_score"`
				Source types.MovieInfo `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
)

func (s *Elasticsearch) SearchMovieByName(ctx context.Context, movieName string) ([]types.MovieInfo, error) {

	// var buf bytes.Buffer
	// sq := SearchQuery{
	// 	Query{
	// 		Match{},
	// 	},
	// }
	// if err := json.NewEncoder(&buf).Encode(sq); err != nil {
	// 	log.Fatalf("Error encoding query: %s", err)
	// }
	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex("movies"),
		s.client.Search.WithQuery("name:"+movieName),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithSize(10),

		s.client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	resEs := ResponseElasticsearch{}
	if err := json.NewDecoder(res.Body).Decode(&resEs); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	if len(resEs.Hits.Hits) == 0 {
		return nil, ErrMovieNotFound
	}
	var list []types.MovieInfo
	for _, h := range resEs.Hits.Hits {
		m := h.Source
		m.ID = h.ID
		list = append(list, m)
	}
	return list, nil
}
