package crawler

import (
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
		// Create(ctx context.Context) (string, error)
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

func NewServiceCLI() *Service {
	return &Service{}
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
		fmt.Println("Visiting", r.URL.String())
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
			crawlMovieInfo(movieLink)
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

func crawlMovieInfo(link string) {
	c := colly.NewCollector()
	// cVideos := c.Clone()
	c.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		title := e.DOM.Find(".title_wrapper h1")
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
				// writers = writers[:len(writers)-1]
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
			// c.Visit(domain + href)

			// c.OnRequest(func(r *colly.Request) {
			// 	fmt.Println("Visiting", r.URL.String())
			// })

			// c.OnHTML(`#imdb-jw-video-1 video[class="jw-video jw-reset"]`, func(e *colly.HTMLElement) {
			// 	src := e.Attr("src")

			// trailerPaths = append(trailerPaths, src)
			// })
			trailerPaths = append(trailerPaths, getTrailerPaths(href))

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
		fmt.Printf("movie: %+v\n", movie)
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

func getTrailerPaths(link string) string {
	c := colly.NewCollector()
	var trailerPaths string
	c.OnHTML(`video[class="jw-video jw-reset"`, func(e *colly.HTMLElement) {

		trailerPaths = e.Attr("src")
		// trailerPaths, _ = e.DOM.Find(`video[class="jw-video jw-reset"]`).Attr("src")

	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Headers.Set("accept-language", "en-US,en;q=0.9")
	})
	c.Visit(domain + link)

	return trailerPaths
}
