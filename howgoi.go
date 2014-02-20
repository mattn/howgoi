package howgoi

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

type Answer struct {
	Code string
	Link string
	Tags []string
}

func getAnswers(u string) (ans *Answer, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	uri, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	uri.Query().Set("answertab", "votes")
	doc, err := goquery.NewDocument(uri.String())
	if err != nil {
		return nil, err
	}
	answer := doc.Find(".answer").First()
	code := answer.Find("pre")
	if code == nil {
		code = answer.Find("code").First()
	}
	if code == nil {
		code = answer.Find(".post-text").First()
	}
	text := code.Text()
	if text == "" {
		return nil, errors.New("Anseer Not Found")
	}
	tags := doc.Find(".post-tag").Map(func(_ int, a *goquery.Selection) string {
		return a.Text()
	})
	return &Answer{
		Code: text,
		Link: u,
		Tags: tags,
	}, nil
}

func Query(args ...string) ([]Answer, error) {
	return QueryN(1, args...)
}

func QueryN(n int, args ...string) (answers []Answer, err error) {
	if n == -1 {
		n = 1
	}
	query := strings.Replace(strings.Join(args, " "), "?", "", -1)

	q := url.Values{}
	q.Set("q", fmt.Sprintf("site:%s %s", "stackoverflow.com", query))
	u := "https://www.google.com/search?" + q.Encode()
	doc, err := goquery.NewDocument(u)
	if err != nil {
		return nil, err
	}
	urls := []string{}
	doc.Find(".r a").Each(func(_ int, a *goquery.Selection) {
		if len(urls) >= n {
			return
		}
		if link, ok := a.Attr("href"); ok {
			if strings.HasPrefix(link, "/") {
				if uri, err := url.ParseRequestURI(link); err == nil {
					link = uri.Query().Get("q")
				}
			}
			urls = append(urls, link)
		}
	})
	if len(urls) == 0 {
		doc.Find(".l a").Each(func(_ int, a *goquery.Selection) {
			if len(urls) >= n {
				return
			}
			if link, ok := a.Attr("href"); ok {
				if strings.HasPrefix(link, "/") {
					if uri, err := url.ParseRequestURI(link); err == nil {
						link = uri.Query().Get("q")
					}
				}
				urls = append(urls, link)
			}
		})
	}

	for _, u = range urls {
		answer, err := getAnswers(u)
		if err == nil {
			answers = append(answers, *answer)
		}
	}
	return answers, nil
}
