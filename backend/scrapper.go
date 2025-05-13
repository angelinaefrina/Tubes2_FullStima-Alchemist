package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type ElementFromFandom struct {
	Name           string     `json:"element"`
	Tier           int        `json:"tier"`
	LocalSVGPath   string     `json:"local_svg_path"`
	OriginalSVGURL string     `json:"original_svg_url"`
	Recipes        [][]string `json:"recipes"`
}

const baseURL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

func ScrapeAll() ([]ElementFromFandom, error) {
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

	tierRegex := regexp.MustCompile(`Tier (\d+)`)

	doc.Find("h3").Each(func(_ int, hdr *goquery.Selection) {
		title := hdr.Find("span.mw-headline").Text()
		if title == "" {
			return
		}

		// Extract tier or fallback to use raw title
		tier := 0
		tierFolder := strings.ReplaceAll(title, " ", "_")
		if matches := tierRegex.FindStringSubmatch(title); len(matches) > 1 {
			tier, _ = strconv.Atoi(matches[1])
			tierFolder = "Tier_" + matches[1] + "_elements"
		} else if title == "Starting Elements" {
			tierFolder = "Starting_elements"
		} else {
			tierFolder = strings.ReplaceAll(title, " ", "_")
		}

		// Locate the associated table
		tbl := hdr.Next()
		for tbl.Length() > 0 && !tbl.Is("table.list-table") {
			tbl = tbl.Next()
		}
		if tbl.Length() == 0 {
			return
		}

		dir := filepath.Join("public", "svgs", tierFolder)
		os.MkdirAll(dir, 0755)

		tbl.Find("tr").Each(func(j int, row *goquery.Selection) {
			if j == 0 {
				return
			}

			wg.Add(1)
			go func(row *goquery.Selection, tier int, tierFolder string) {
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
				imgNode := cols.Eq(0).Find("img")

				// Coba ambil src
				imgSrc, exists := imgNode.Attr("src")

				// Jika src base64 atau kosong, coba fallback ke data-src
				if !exists || strings.HasPrefix(imgSrc, "data:") || imgSrc == "" {
					imgSrc, exists = imgNode.Attr("data-src")
				}

				// Jika masih tidak valid, coba ambil dari srcset
				if !exists || imgSrc == "" {
					if srcset, ok := imgNode.Attr("srcset"); ok && srcset != "" {
						imgSrc = strings.Split(srcset, " ")[0] // ambil yang pertama dari daftar
						exists = true
					}
				}

				// Pastikan imgSrc valid, lalu buat URL lengkap
				if exists && imgSrc != "" && !strings.HasPrefix(imgSrc, "data:") {
					if strings.HasPrefix(imgSrc, "//") {
						imgURL = "https:" + imgSrc
					} else if strings.HasPrefix(imgSrc, "/") {
						imgURL = "https://little-alchemy.fandom.com" + imgSrc
					} else {
						imgURL = imgSrc
					}
				}

				// Sanitize filename
				safeName := strings.NewReplacer(" ", "_", "/", "_", "\\", "_", ":", "_", "*", "_", "?", "_", "\"", "_", "<", "_", ">", "_", "|", "_").Replace(name)
				localPath := filepath.Join(tierFolder, safeName+".svg")
				fullPath := filepath.Join("public", "svgs", localPath)

				// Download image if valid
				if imgURL != "" {
					downloadSVG(imgURL, fullPath)
				}

				// Extract recipes
				recipes := [][]string{}
				cols.Eq(1).Find("ul li").Each(func(_ int, li *goquery.Selection) {
					parts := li.Find("a[title]").Map(func(_ int, a *goquery.Selection) string {
						return a.Text()
					})
					if len(parts) == 2 {
						recipes = append(recipes, parts)
					}
				})

				// Append element
				mu.Lock()
				allElements = append(allElements, ElementFromFandom{
					Name:           name,
					Tier:           tier,
					LocalSVGPath:   localPath,
					OriginalSVGURL: imgURL,
					Recipes:        recipes,
				})
				mu.Unlock()

			}(row, tier, tierFolder)
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
		println("ERROR: Failed to GET", url, "-", err.Error())
		return
	}
	defer resp.Body.Close()

	os.MkdirAll(filepath.Dir(dest), 0755)
	f, err := os.Create(dest)
	if err != nil {
		println("ERROR: Failed to CREATE FILE", dest, "-", err.Error())
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		println("ERROR: Failed to WRITE FILE", dest, "-", err.Error())
	}
}
