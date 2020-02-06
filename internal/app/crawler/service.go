package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dqkcode/movie-database/internal/app/types"
	"github.com/google/uuid"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type (
	movieService interface {
		CreateMovie(ctx context.Context, movie types.MovieInfo) error
	}
	Service struct {
		movieService
	}
)

const domain = "https://www.imdb.com"

func NewService(movieService movieService) *Service {
	return &Service{
		movieService,
	}
}

func NewServiceCLI(movieService movieService) *Service {
	return &Service{
		movieService,
	}
}

func (s *Service) GetAllGenres() []string {
	c := colly.NewCollector(
	// colly.AllowedDomains("imdb.com"),
	)
	var links []string
	c.OnHTML("a[name=slot_center-6]", func(e *colly.HTMLElement) {

		e.DOM.Next().Find(".table-row").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Find("a").Attr("href")
			links = append(links, href)
		})

	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(domain + "/feature/genre/")
	return links
}

func (s *Service) CrawlAllMovies(ctx context.Context) ([]*types.MovieInfo, error) {
	links := s.GetAllGenres()
	var movies []*types.MovieInfo
	c := colly.NewCollector()

	for _, link := range links {
		c.OnHTML("h3.lister-item-header", func(e *colly.HTMLElement) {
			movieLink := e.ChildAttr("a", "href")
			s.CrawlMovieInfo(ctx, movieLink)

		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		c.Visit(domain + link)
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	return movies, nil
}

func (s *Service) CrawlMovieInfo(ctx context.Context, link string) {
	c := colly.NewCollector()
	c.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		title := e.DOM.Find(".title_wrapper h1")
		name := e.DOM.Find(".title_wrapper h1").Text()
		fmt.Printf("Crawling movie: %s\n", name)
		subText := title.Next()
		movieLength := subText.Find("time[datetime]").Text()
		var genres, casts, writers, images, trailerPaths []string
		e.DOM.Find(`#titleStoryLine div[class="see-more inline canwrap"]`).Last().Find("a").Each(func(_ int, s *goquery.Selection) {
			genres = append(genres, s.Text())

		})
		releaseTime := e.DOM.Find(`a[title="See more release dates"]`).Text()

		e.DOM.Find(".plot_summary .credit_summary_item").Eq(1).Each(func(_ int, s *goquery.Selection) {
			s.Find("a").Each(func(_ int, s *goquery.Selection) {
				writers = append(writers, s.Text())
			})
		})

		rate := e.DOM.Find(`.ratingValue span[itemprop="ratingValue"]`).Text()

		e.DOM.Find(".cast_list tbody tr").Each(func(_ int, s *goquery.Selection) {
			star := s.Find("td:nth-child(2) a").Text()
			star = strings.TrimSpace(star)
			casts = append(casts, star)
		})

		e.DOM.Find(`div[class="mediastrip_big"] span a`).Each(func(_ int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			trailerPaths = append(trailerPaths, getTrailerPaths(c, href))

		})

		e.DOM.Find(`div[class="mediastrip"] a img`).Each(func(_ int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			images = append(images, src)
		})
		storyLine := strings.TrimSpace(e.DOM.Find(`#titleStoryLine div[class="inline canwrap"]`).First().Find("p span").Text())

		movie := types.MovieInfo{
			ID:           uuid.New().String(),
			Name:         e.DOM.Find(".title_wrapper h1").Text(),
			MovieLength:  convertTimeToInt(movieLength),
			ReleaseTime:  releaseTime,
			Director:     e.DOM.Find(".plot_summary .credit_summary_item a").First().Text(),
			Writers:      writers,
			Rate:         rate,
			Genres:       genres,
			Casts:        casts[1:len(casts)],
			Storyline:    storyLine,
			ImagesPath:   images,
			TrailersPath: trailerPaths,
		}
		NewUser := &types.UserInfo{
			ID: "123456789",
			//TODO add some info
		}
		newCtx := context.WithValue(ctx, "crawler", NewUser)

		err := s.movieService.CreateMovie(newCtx, movie)
		if err != nil {
			fmt.Printf("Err save movie %s:  %s\n", name, err)
		}
		fmt.Printf("Crawled movie :%v", name)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Headers.Set("accept-language", "en-US,en;q=0.9")
	})

	c.Visit(domain + link)
}

func convertTimeToInt(time string) int {
	time = strings.TrimSpace(time)
	t := strings.Split(time, "h ")
	hh, err := strconv.Atoi(t[0])
	if err != nil {
		return -1
	}
	mm, err := strconv.Atoi(t[1][:2])
	if err != nil {
		return -1
	}
	result := hh*60 + mm

	return result
}

type (
	Response struct {
		VideoPlaybackId      string `json:"videoPlaybackId"`
		VideoPlaybackToken   string `json:"videoPlaybackToken"`
		VideoLegacyAdUrl     string `json:"videoLegacyAdUrl"`
		VideoLegacyEncodings []struct {
			Definition string `json:"definition"`
			MimeType   string `json:"mimeType"`
			URL        string `json:"url"`
		} `json:"videoLegacyEncodings"`
	}
)

func getTrailerPaths(c *colly.Collector, link string) string {
	// c := colly.NewCollector()
	var trailerPath string
	c.OnHTML(`#a-page`, func(e *colly.HTMLElement) {
		v := c.Clone()
		scriptText := e.DOM.Find("script").Text()
		videoInfoKey := scriptText[strings.Index(scriptText, `playbackDataKey":["`)+19 : strings.Index(scriptText, `"],"activeWeblabs`)]
		v.OnResponse(func(r *colly.Response) {
			res := []Response{}
			if err := json.Unmarshal(r.Body, &res); err != nil {
				fmt.Println("Error: ", err)
			}
			if len(res[0].VideoLegacyEncodings) != 0 {
				trailerPath = res[0].VideoLegacyEncodings[len(res[0].VideoLegacyEncodings)-1].URL
			}

		})

		v.Visit(domain + "/ve/data/VIDEO_PLAYBACK_DATA?key=" + videoInfoKey)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("accept-language", "en-US,en;q=0.9")
	})
	c.Visit(domain + link)

	return trailerPath
}
