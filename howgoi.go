package howgoi

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func Query(args... string) ([]string, error) {
	return QueryN(1, args...)
}

func QueryN(n int, args... string) ([]string, error) {
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

	codes := []string{}
	for _, u = range urls {
		uri, err := url.Parse(u)
		if err != nil {
			continue
		}
		uri.Query().Set("answertab", "votes")
		doc, err = goquery.NewDocument(uri.String())
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
		codes = append(codes, code.Text())
	}
	return codes, nil
}
