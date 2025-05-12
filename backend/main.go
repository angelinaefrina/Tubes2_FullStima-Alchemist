package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type TreeResponse struct {
	Trees        []*Node `json:"trees"`
	NodesVisited int     `json:"nodesVisited"`
	SearchTime   float64 `json:"searchTime"` // detik
}

var recipeData []ElementFromFandom

func main() {
	// Jalankan scraper dulu
	log.Println("Scraping data...")
	data, err := ScrapeAll()
	if err != nil {
		log.Fatalf("scraping failed: %v", err)
	}
	recipeData = data
	log.Printf("Scraped %d elements\n", len(recipeData))
	// fmt.Println(recipeData)

	http.HandleFunc("/api/recipe", recipeHandler)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type SearchRequest struct {
	Method     string `json:"method"`     // "dfs" atau "bfs"
	Target     string `json:"target"`     // nama elemen
	Multiple   bool   `json:"multiple"`   // true kalau cari banyak
	MaxRecipes int    `json:"maxRecipes"` // jumlah maksimal resep (opsional)
}

func recipeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request method: %s", r.Method)

	// Allow CORS for all origins
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight (OPTIONS) request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decode JSON request body
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Searching for element: %s", req.Target)

	target := strings.TrimSpace(req.Target)
	if target == "" {
		http.Error(w, "missing target", http.StatusBadRequest)
		return
	}

	// Cek apakah elemen target ada di dataset
	found := false
	for _, el := range recipeData {
		if el.Name == target {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "element not found in data", http.StatusNotFound)
		return
	}

	// Build peta resep
	log.Println("Building recipe map...")
	recipeMap := buildRecipeMap(recipeData)
	tierNum := buildTierMap(recipeData)

	// Identifikasi base elements (tidak punya resep)
	baseElements := map[string]bool{}
	for _, el := range recipeData {
		//if len(el.Recipes) == 0 {
		if el.Name == "Air" || el.Name == "Earth" || el.Name == "Fire" || el.Name == "Water" {
			baseElements[el.Name] = true
		}
	}

	fmt.Printf("Base elements: %d\n", len(baseElements))
	fmt.Printf("Recipes: %d\n", len(recipeMap))
	fmt.Printf("Target: %s\n", target)
	var trees []*Node
	var totalNodes int
	var err error

	start := time.Now() // mulai hitung waktu

	if req.Multiple {
		var paths interface{}
		if req.Method == "bfs" {
			log.Print("Finding multiple paths using BFS...")
			var bfsPaths map[string][][]string
			bfsPaths, err = bfsMultiplePaths(target, recipeMap, baseElements, req.MaxRecipes, tierNum)
			paths = bfsPaths
		} else if req.Method == "dfs" {
			log.Print("Finding multiple paths using DFS...")
			var dfsPaths [][]string
			dfsPaths, err = dfsMultiplePaths(target, recipeMap, baseElements, req.MaxRecipes, tierNum)
			paths = dfsPaths
		} else {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}

		// log.Printf("Found %v paths", paths)
		switch p := paths.(type) {
		case map[string][][]string:
			for _, pathList := range p[target] {
				tree := pathToTree(pathList)
				trees = append(trees, tree)
				totalNodes += int(countNodes(tree))
			}
		case [][]string:
			for _, path := range p {
				tree := pathToTree(path)
				trees = append(trees, tree)
				totalNodes += int(countNodes(tree))
			}
		}
	} else {
		var path []string
		if req.Method == "bfs" {
			log.Print("Finding single path using BFS...")
			path, err = bfs(target, recipeMap, baseElements, map[string]int{})
		} else if req.Method == "dfs" {
			log.Print("Finding single path using DFS...")
			path, err = dfs(target, recipeMap, baseElements, tierNum)
		} else {
			http.Error(w, "invalid method", http.StatusBadRequest)
			return
		}
		tree := pathToTree(path)
		trees = []*Node{tree}
		totalNodes = int(countNodes(tree))
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	elapsed := time.Since(start).Seconds()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TreeResponse{
		Trees:        trees,
		NodesVisited: totalNodes,
		SearchTime:   elapsed,
	})
}
