package main

//            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
//                    Version 2, December 2004
//
// Copyright (C) 2024 Thiago Kenji Okada <thiagokokada@gmail.com>
//
// Everyone is permitted to copy and distribute verbatim or modified
// copies of this license document, and changing it is allowed as long
// as the name is changed.
//
//            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
//   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION
//
//  0. You just DO WHAT THE FUCK YOU WANT TO.

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	"github.com/gosimple/slug"
	"github.com/yuin/goldmark"
)

const rssBaseUrl = "https://github.com/thiagokokada/blog/blob/main"
const readmeTemplate = `# Blog

Mirror of my blog in https://kokada.capivaras.dev/.

## Posts

[![RSS](https://img.shields.io/badge/RSS-FFA562?style=for-the-badge&logo=rss&logoColor=white)](https://raw.githubusercontent.com/thiagokokada/blog/main/rss.xml)

%s
`

type post struct {
	title    string
	file     string
	slug     string
	contents []byte
	date     time.Time
}

func must1[T any](v T, err error) T {
	must(err)
	return v
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func extractTitleAndContents(raw []byte) (title string, contents []byte, err error) {
	for i, c := range raw {
		// We are assuming that each file has one title as a H1 header
		if c == '\n' {
			original := string(raw[:i])
			clean := strings.Replace(string(raw[:i]), "# ", "", 1)
			if clean == original {
				return "", nil, fmt.Errorf("could not find title")
			}
			title = clean
			contents = bytes.TrimSpace(raw[i:])
			break
		}
	}

	return title, contents, nil
}

func grabPosts() []post {
	var posts []post
	must(filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignore hidden files
		if d.Name()[0] == '.' {
			return nil
		}
		// Ignore non-Markdown files
		if filepath.Ext(d.Name()) != ".md" {
			return nil
		}

		// Get the directory of the file
		dir := filepath.Dir(path)
		// Ignore files in the current directory
		if dir == "." {
			return nil
		}

		// Parse directory name as a date
		date, err := time.Parse(time.DateOnly, dir)
		if err != nil {
			log.Printf("[WARN]: ignoring non-date directory: %s\n", path)
			return nil
		}

		// Load the contents of the Markdown and try to parse
		// the title
		raw := must1(os.ReadFile(path))
		title, contents, err := extractTitleAndContents(raw)
		if err != nil || title == "" || contents == nil {
			return fmt.Errorf(
				"something is wrong with file: %s, title: %s, contents: %s, error: %w",
				path,
				title,
				contents,
				err,
			)
		}

		posts = append(posts, post{
			title:    title,
			file:     path,
			slug:     slug.Make(title),
			contents: contents,
			date:     date,
		})

		return nil
	}))

	sort.Slice(posts, func(i, j int) bool {
		// Reverse sorting based on filename
		return posts[i].file > posts[j].file
	})

	return posts
}

func genRss(posts []post) string {
	feed := &feeds.Feed{
		Title:       "kokada's blog",
		Description: "# dd if=/dev/urandom of=/dev/brain0",
	}
	var items []*feeds.Item
	for _, post := range posts {
		link := must1(url.JoinPath(rssBaseUrl, post.file))
		var buf bytes.Buffer
		must(goldmark.Convert(post.contents, &buf))
		items = append(items, &feeds.Item{
			Title:       post.title,
			Link:        &feeds.Link{Href: link},
			Created:     post.date,
			Id:          link,
			Description: buf.String(),
		})
	}
	feed.Items = items
	return must1(feed.ToRss())
}

func genReadme(posts []post) string {
	var titles []string
	for _, post := range posts {
		title := fmt.Sprintf(
			"- [%s](%s) - %s",
			post.title,
			post.file,
			post.date.Format(time.DateOnly),
		)
		titles = append(titles, title)
	}
	return fmt.Sprintf(readmeTemplate, strings.Join(titles, "\n"))
}

func main() {
	rss := flag.Bool("rss", false, "Generate RSS (XML) instead of README.md")
	publish := flag.Bool("publish", false, "Publish updates to Maratoa instance")
	flag.Parse()

	posts := grabPosts()
	if *publish {
		if mataroaToken == "" {
			log.Println("[WARN]: empty MATAROA_TOKEN env var")
		}
		publishToMataroa(posts)
	} else if *rss {
		fmt.Print(genRss(posts))
	} else {
		fmt.Print(genReadme(posts))
	}
}
