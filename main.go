package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Papeis []string `yaml:"papeis"`
}

var weirdPrefix = regexp.MustCompile(`^[^\p{L}\p{N}]+`)

func main() {
	config := loadConfig("papeis.yaml")

	for _, papel := range config.Papeis {
		dados, err := scrapePapel(papel)
		if err != nil {
			log.Printf("erro ao buscar %s: %v\n", papel, err)
			continue
		}

		fmt.Println("==========")
		fmt.Println("Papel:", papel)
		fmt.Println("Cotação:", dados["Cotação"])
		fmt.Println("P/L:", dados["P/L"])
		fmt.Println("P/VP:", dados["P/VP"])
		fmt.Println("Div. Yield:", dados["Div. Yield"])
		fmt.Println("ROE:", dados["ROE"])
		fmt.Println("Setor:", dados["Setor"])
	}
}

func loadConfig(path string) Config {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func scrapePapel(papel string) (map[string]string, error) {
	url := "https://fundamentus.com.br/detalhes.php?papel=" + papel

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	utf8Reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(utf8Reader)
	if err != nil {
		return nil, err
	}

	dados := map[string]string{}

	doc.Find("td.label").Each(func(i int, s *goquery.Selection) {
		label := clean(s.Text())
		value := clean(s.Next().Text())

		dados[label] = value
	})

	return dados, nil
}

func clean(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	s = weirdPrefix.ReplaceAllString(s, "")

	return s
}
