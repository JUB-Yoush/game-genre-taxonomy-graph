package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"slices"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"

	// "github.com/yuin/goldmark/extension"
	meta "github.com/yuin/goldmark-meta"
)

type Genre struct {
	id        string
	parents   []string
	children  []string
	mdContent bytes.Buffer
}

type Game struct {
	title   string
	release time.Time
	link    string
}

func newGenre(id string, content string) {
}

func main() {
	genreMap := map[string]*Genre{}
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	rootPath := "public"
	genreDataFolder := "markdown_data"
	os.RemoveAll(rootPath)

	// read in all files from posts directory
	files, err := os.ReadDir(genreDataFolder)
	if err != nil {
		log.Fatal(err)
	}

	// create genre objects
	for _, f := range files {
		fmt.Println(f.Name())
		context := parser.NewContext()
		byteArr, err := os.ReadFile(fmt.Sprintf("%s/%s", genreDataFolder, f.Name()))

		var buf bytes.Buffer

		if err != nil {
			log.Fatal(err)
		}

		if err := markdown.Convert([]byte(byteArr[:]), &buf, parser.WithContext(context)); err != nil {
			panic(err)
		}

		metaData := meta.Get(context)

		title := metaData["id"]
		parents := metaData["parents"]
		genre := new(Genre)
		genre.id = title.(string)
		genre.mdContent = buf

		for _, v := range parents.([]any) {
			s := v.(string)
			genre.parents = append(genre.parents, s)
		}
		genreMap[genre.id] = genre
	}

	// map parent genres to child genres
	for id, genre := range genreMap {
		parents := genre.parents
		for id2, genre2 := range genreMap {
			if slices.Contains(parents, id2) {
				genre2.children = append(genre2.children, id)
			}
		}
	}
}
