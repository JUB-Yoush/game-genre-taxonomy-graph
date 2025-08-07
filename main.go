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
	id         string
	name       string
	parents    []string
	classNames string
	children   []*Genre
	rawMd      bytes.Buffer
	html       templ.Component
}

func (g Genre) String() string {
	return fmt.Sprint(g.id)
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

func makeGenreList() (map[string]*Genre, []string) {
	genreMap := map[string]*Genre{}
	var rootGenres []string
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
		fmt.Println("-----")
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
		genre.rawMd = buf
		genre.name = name.(string)
		genre.html = Unsafe(genre.rawMd.String())

		genre.classNames = ""

		if parents == nil {
			fmt.Println("no parents")
			rootGenres = append(rootGenres, genre.id)
			genreMap[genre.id] = genre
			continue
		}

		for _, v := range parents.([]any) {
			s := v.(string)
			genre.parents = append(genre.parents, s)
			genre.classNames += fmt.Sprintf("%s ", s)
		}

		fmt.Println(genre.classNames)
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
	return genreMap, rootGenres
}

func makeGenreTree(genreMap map[string]*Genre, rootGenres []string) (res [][]*Genre) {
	var queue []*Genre
	visited := map[*Genre]bool{}

	for _, genre := range rootGenres {
		queue = append(queue, genreMap[genre])
	}
	res = append(res, queue)

	var row []*Genre
	for len(queue) > 0 {
		for _, curr := range queue {
			_, in := visited[curr]
			if in {
				continue
			}
			visited[curr] = true
			row = append(row, curr.children...)
		}
		res = append(res, row)
		queue = row
		row = nil
	}
	res = res[:len(res)-1]
	return res
}

func main() {
	genreMap, rootGenres := makeGenreList()
	genreTree := makeGenreTree(genreMap, rootGenres)
	fmt.Println(genreTree)
	component := BoilerPlate(genreMap, genreTree)

	rootPath := "public"
	staticPath := "static"
	os.RemoveAll(rootPath)
	if err := os.Mkdir(rootPath, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	err := os.CopyFS(rootPath, os.DirFS(staticPath))
	if err != nil {
		log.Fatal(err)
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
