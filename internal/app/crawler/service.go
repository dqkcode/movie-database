package crawler

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/dqkcode/movie-database/internal/app/types"

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

	c.OnHTML("#content-2-wide", func(e *colly.HTMLElement) {
		title := e.DOM.Find(".title_wrapper h1")
		subText := title.Next()
		time := subText.Find("time[datetime]").Text()
		convertTime(time)
		movie := types.MovieInfo{
			ID:   uuid.New().String(),
			Name: e.DOM.Find(".title_wrapper h1").Text(),
		}
		fmt.Println(movie)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
		r.Headers.Set("accept-language", "en-US,en;q=0.9")
	})

	c.Visit(domain + link)
}

// 2h 15min

func convertTime(time string) int {
	time = strings.Trim(time, " ")
	t := strings.Split(time, "h ")
	fmt.Println(t)
	return 1
}
