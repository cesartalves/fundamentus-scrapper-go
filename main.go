package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"

	"regexp"
)

var leadingQuestionMark = regexp.MustCompile(`^\?+`)

func main() {
	papel := "TAEE11"

	url := "https://fundamentus.com.br/detalhes.php?papel=" + papel

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Corrige encoding ISO-8859-1 -> UTF-8
	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		log.Fatal(err)
	}

	dados := map[string]string{}

	doc.Find("td.label").Each(func(i int, s *goquery.Selection) {
		label := clean(s.Text())

		value := clean(s.Next().Text())

		dados[label] = value
	})

	for k, v := range dados {
		fmt.Printf("%s => %s\n", k, v)
	}
}

func clean(s string) string {
	s = strings.TrimSpace(s)

	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")

	// remove ? no começo
	s = leadingQuestionMark.ReplaceAllString(s, "")

	return s
}
