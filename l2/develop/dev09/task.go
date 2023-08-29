package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

/*
Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

/*
	Команда запуска имеет следующий вид:
	task [OPTION...] url
	url - адрес начальной страницы

	OPTION:
		-r - рекурсивно проходить по ссылкам <a>,
		-rps <int> - ограничить количество запросов в секунду (default = 1000)
*/

type Config struct {
	Url *url.URL
	Rps int
	Rec bool
}

func parseConfig() *Config {

	var cfg Config
	var err error

	flag.BoolVar(&cfg.Rec, "r", false, "рекурсивно проходить по ссылкам <a>")
	flag.IntVar(&cfg.Rps, "rps", 1000, "ограничить количество запросов в секунду (default = 1000)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "can't find url")
		flag.Usage()
		os.Exit(2)
	}

	cfg.Url, err = url.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't parse url: %s\n", err)
		flag.Usage()
		os.Exit(2)
	}

	return &cfg
}

func savefile(p string, b []byte) error {
	dir, file := path.Split(p)
	if file == "" {
		file = "index.html"
	}
	// Создаем путь до файла
	if err := os.MkdirAll(dir, os.FileMode(0644)); err != nil {
		return err
	}
	// Сохраняем файл
	p = path.Join(dir, file)
	if err := os.WriteFile(p, b, 0644); err != nil {
		return err
	}
	return nil
}

func wget(cfg *Config, logger *log.Logger) {
	c := colly.NewCollector(
		colly.AllowedDomains(cfg.Url.Host),
		colly.Async(true),
	)

	err := c.Limit(
		&colly.LimitRule{
			Delay: time.Duration(int(time.Second) / cfg.Rps),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	if cfg.Rec {
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			trimedHref := strings.TrimSpace(e.Attr("href"))
			sep := strings.SplitN(trimedHref, "#", 2)
			if len(sep) == 0 || len(sep[0]) == 0 {
				return
			}
			link := e.Request.AbsoluteURL(sep[0])
			if err := c.Visit(link); err != nil {
				skip := errors.Is(err, colly.ErrAlreadyVisited) ||
					errors.Is(err, colly.ErrForbiddenDomain)

				if !skip {
					logger.Println(err, link)
				}
			}
		})
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		u, err := url.JoinPath(r.Request.URL.Host, r.Request.URL.Path)
		if err != nil {
			logger.Println(err)
			return
		}
		if err := savefile(u, r.Body); err != nil {
			logger.Println(err)
		}
	})

	if err := c.Visit(cfg.Url.String()); err != nil {
		logger.Println(err)
	}
	c.Wait()
}

func main() {
	cfg := parseConfig()
	logger := log.New(os.Stderr, "error: ", log.Ldate|log.Ltime|log.Lshortfile)
	wget(cfg, logger)
}
