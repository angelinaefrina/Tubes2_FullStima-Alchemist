package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	// "strings"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

type ElementFromFandom struct {
	Name           string     `json:"name"`
	LocalSVGPath   string     `json:"local_svg_path"`   // will be empty
	OriginalSVGURL string     `json:"original_svg_url"` // URL only, not downloaded
	Recipes        [][]string `json:"recipes"`
}

func scraper() {
	// 1) Scrape & write JSON
	data, err := scrapeAll()
	if err != nil {
		log.Fatalf("scrape failed: %v", err)
	}

	raw, _ := json.MarshalIndent(data, "", "  ")
	if err := os.WriteFile("recipe.json", raw, 0644); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote recipe.json (%d sections)", len(data))

	// 2) HTTP handlers
	http.HandleFunc("/api/recipes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(raw)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func scrapeAll() (map[string][]ElementFromFandom, error) {
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	out := make(map[string][]ElementFromFandom)
	doc.Find("h3").Each(func(_ int, hdr *goquery.Selection) {
		title := hdr.Find("span.mw-headline").Text()
		if title == "" {
			return
		}

		tbl := hdr.Next()
		for tbl.Length() > 0 && !tbl.Is("table.list-table") {
			tbl = tbl.Next()
		}
		if tbl.Length() == 0 {
			return
		}

		sectionKey := title
		var elems []ElementFromFandom

		tbl.Find("tr").Each(func(i int, row *goquery.Selection) {
			if i == 0 {
				return
			}
			cols := row.Find("td")
			if cols.Length() < 2 {
				return
			}
			name := cols.Eq(0).Find("a[title]").First().Text()
			if name == "" {
				return
			}

			// Only get the original SVG URL (optional)
			fileA := cols.Eq(0).Find("a.mw-file-description")
			href, _ := fileA.Attr("href")

			recipes := [][]string{}
			cols.Eq(1).Find("ul li").Each(func(_ int, li *goquery.Selection) {
				parts := li.Find("a[title]").Map(func(_ int, a *goquery.Selection) string {
					return a.Text()
				})
				if len(parts) == 2 {
					recipes = append(recipes, []string{parts[0], parts[1]})
				}
			})

			elems = append(elems, ElementFromFandom{
				Name:           name,
				LocalSVGPath:   "",
				OriginalSVGURL: href,
				Recipes:        recipes,
			})
		})

		if len(elems) > 0 {
			out[sectionKey] = elems
		}
	})

	return out, nil
}
