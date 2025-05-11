package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type ElementFromFandom struct {
	Name           string     `json:"element"`
	LocalSVGPath   string     `json:"local_svg_path"`
	OriginalSVGURL string     `json:"original_svg_url"`
	Recipes        [][]string `json:"recipes"`
}

const baseURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

func ScrapeAll() ([]ElementFromFandom, error) {
	// Gunakan custom User-Agent agar tidak diblokir
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allElements []ElementFromFandom

	doc.Find("table.list-table").Each(func(i int, tbl *goquery.Selection) {
		tbl.Find("tr").Each(func(j int, row *goquery.Selection) {
			if j == 0 {
				return
			}
			wg.Add(1)
			go func(row *goquery.Selection) {
				defer wg.Done()

				cols := row.Find("td")
				if cols.Length() < 2 {
					return
				}

				name := cols.Eq(0).Find("a[title]").First().Text()
				if name == "" {
					return
				}

				var imgURL string
				imgSrc, exists := cols.Eq(0).Find("img").Attr("src")
				if exists {
					if strings.HasPrefix(imgSrc, "//") {
						imgURL = "https:" + imgSrc
					} else if strings.HasPrefix(imgSrc, "/") {
						imgURL = "https://little-alchemy.fandom.com" + imgSrc
					} else {
						imgURL = imgSrc
					}
				}

				localPath := ""
				if imgURL != "" {
					fname := strings.ReplaceAll(name, " ", "_") + filepath.Ext(imgURL)
					localPath = fname
					downloadSVG(imgURL, "svgs/"+fname)
				}

				recipes := [][]string{}
				cols.Eq(1).Find("ul li").Each(func(_ int, li *goquery.Selection) {
					parts := li.Find("a[title]").Map(func(_ int, a *goquery.Selection) string {
						return a.Text()
					})
					if len(parts) == 2 {
						recipes = append(recipes, parts)
					}
				})

				mu.Lock()
				allElements = append(allElements, ElementFromFandom{
					Name:           name,
					LocalSVGPath:   localPath,
					OriginalSVGURL: imgURL,
					Recipes:        recipes,
				})
				mu.Unlock()
			}(row)
		})
	})

	wg.Wait()
	return allElements, nil
}

func downloadSVG(url, dest string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	os.MkdirAll(filepath.Dir(dest), 0755)
	f, err := os.Create(dest)
	if err != nil {
		return
	}
	defer f.Close()

	io.Copy(f, resp.Body)
}
