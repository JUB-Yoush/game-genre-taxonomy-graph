package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"time"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"

	// "github.com/yuin/goldmark/extension"
	meta "github.com/yuin/goldmark-meta"
)

type Genre struct {
	id        string
	name      string
	parents   []string
	children  []*Genre
	mdContent bytes.Buffer
}

type Game struct {
	title   string
	release time.Time
	link    string
}

func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}

func makeGenreList() map[string]*Genre {
	genreMap := map[string]*Genre{}
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	genreDataFolder := "markdown_data"

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
		name := metaData["name"]
		genre := new(Genre)
		genre.id = title.(string)
		genre.mdContent = buf
		genre.name = name.(string)

		if parents == nil {
			genreMap[genre.id] = genre
			continue
		}

		for _, v := range parents.([]any) {
			s := v.(string)
			genre.parents = append(genre.parents, s)
		}

		genreMap[genre.id] = genre
	}

	// make genre tree of parents and children
	for _, genre := range genreMap {
		parents := genre.parents
		for id2, genre2 := range genreMap {
			if slices.Contains(parents, id2) {
				genre2.children = append(genre2.children, genre)
			}
		}
	}
	return genreMap
}

func main() {
	genreMap := makeGenreList()
	component := BoilerPlate(genreMap)

	rootPath := "public"
	os.RemoveAll(rootPath)
	if err := os.Mkdir(rootPath, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	// Create an index page.
	name := path.Join(rootPath, "index.html")
	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	// Write it out.
	err = component.Render(context.Background(), f)
	if err != nil {
		log.Fatalf("failed to write index page: %v", err)
	}
}
