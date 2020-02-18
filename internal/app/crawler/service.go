package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dqkcode/movie-database/internal/app/types"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

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

	movieService interface {
		CreateMovie(movie types.MovieInfo) error
		GetMovieByName(name string) (*types.MovieInfo, error)
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

	c.Visit(domain + "/feature/genre/")
	return links
}

func (s *Service) CrawlAllMovies() ([]*types.MovieInfo, error) {
	links := s.GetAllGenres()
	var movies []*types.MovieInfo
	c := colly.NewCollector()

	for _, link := range links {
		c.OnHTML("h3.lister-item-header", func(e *colly.HTMLElement) {
			movieLink := e.ChildAttr("a", "href")
			s.CrawlMovieInfo(movieLink)

		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})

		c.Visit(domain + link)
	}

	return movies, nil
}

func (s *Service) CrawlMovieInfo(link string) {
	c := colly.NewCollector()
	c.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		var genres, casts, writers, images, trailerPaths, directors []string
		title := e.DOM.Find(".title_wrapper h1")
		n := strings.TrimSpace(e.DOM.Find(".title_wrapper h1").Text())
		name := n[:len(n)-8]
		releaseTime := strings.TrimSpace(e.DOM.Find(`a[title="See more release dates"]`).Text())
		mv, err := s.GetMovieByName(name)
		fmt.Printf("mv.Name :%v, name:%v, == %v \n", mv.Name, name, mv.Name == name)
		if errors.Is(err, types.ErrMovieNotFound) || mv.ReleaseTime != releaseTime {

			fmt.Printf("\tCrawling movie ❚❚: %s\n", name)
			subText := title.Next()
			movieLength := subText.Find("time[datetime]").Text()
			e.DOM.Find(`#titleStoryLine div[class="see-more inline canwrap"]`).Last().Find("a").Each(func(_ int, s *goquery.Selection) {
				genres = append(genres, strings.TrimSpace(s.Text()))

			})

			e.DOM.Find(".plot_summary .credit_summary_item").Eq(1).Each(func(_ int, s *goquery.Selection) {
				s.Find("a").Each(func(_ int, s *goquery.Selection) {
					writers = append(writers, s.Text())
				})
			})

			e.DOM.Find(".plot_summary .credit_summary_item").First().Find("a").Each(func(_ int, s *goquery.Selection) {

				directors = append(directors, s.Text())

			})

			rate, err := strconv.ParseFloat(e.DOM.Find(`.ratingValue span[itemprop="ratingValue"]`).Text(), 8)
			if err != nil {
				fmt.Println("err parse rate info")
			}

			e.DOM.Find(".cast_list tbody tr").Each(func(_ int, s *goquery.Selection) {
				star := s.Find("td:nth-child(2) a").Text()
				star = strings.TrimSpace(star)
				casts = append(casts, star)
			})

			e.DOM.Find(`div[class="mediastrip_big"] span a`).Each(func(_ int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				path := getTrailerPaths(c, href)
				if path != "" {
					trailerPaths = append(trailerPaths, path)
				}

			})

			// e.DOM.Find(`div[class="mediastrip"] a`).Each(func(_ int, s *goquery.Selection) {
			// 	src, _ := s.Children().Attr("src")
			// 	images = append(images, src)
			// })

			storyLine := strings.TrimSpace(e.DOM.Find(`#titleStoryLine div[class="inline canwrap"]`).First().Find("p span").Text())

			movie := types.MovieInfo{
				Name:         name,
				MovieLength:  convertTimeToInt(movieLength),
				ReleaseTime:  releaseTime,
				Directors:    directors,
				Writers:      writers,
				Rate:         rate,
				Genres:       genres,
				Casts:        casts[1:len(casts)],
				Storyline:    storyLine,
				ImagesPath:   images,
				TrailersPath: trailerPaths,
			}

			err = s.movieService.CreateMovie(movie)
			if err != nil {
				fmt.Printf("\tErr save movie %s:  %s\n", name, err)
			}
			fmt.Printf("\tCrawled movie ☑ :%v\n", name)
		}

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
	mm, err := strconv.Atoi(t[1][:len(t[1])-3])
	if err != nil {
		return -1
	}
	result := hh*60 + mm

	return result
}

func getTrailerPaths(c *colly.Collector, link string) string {
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
			if len(res) != 0 && len(res[0].VideoLegacyEncodings) != 0 {
				trailerPath = res[0].VideoLegacyEncodings[len(res[0].VideoLegacyEncodings)-1].URL
			}
		})
		v.Visit(domain + "/ve/data/VIDEO_PLAYBACK_DATA?key=" + videoInfoKey)
	})

	c.Visit(domain + link)

	return trailerPath
}
