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
	"slices"
	"strings"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/gorilla/feeds"
	"github.com/gosimple/slug"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

const (
	blogBaseUrl    = "https://github.com/thiagokokada/blog"
	blogMainUrl    = blogBaseUrl + "/blob/main/"
	blogRawUrl     = blogBaseUrl + "/raw/main/"
	highlightStyle = "monokai"
	readmeTemplate = `# Blog

Mirror of my blog in https://kokada.capivaras.dev/.

## Posts

[![RSS](https://img.shields.io/badge/RSS-FFA562?style=for-the-badge&logo=rss&logoColor=white)](https://raw.githubusercontent.com/thiagokokada/blog/main/rss.xml)

%s
`
)

type post struct {
	title    string
	slug     string
	contents []byte
	date     time.Time
}
type path = string
type posts = *orderedmap.OrderedMap[path, post]

func must2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	must(err)
	return v1, v2
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
	if len(raw) == 0 {
		return "", nil, fmt.Errorf("empty file")
	}

	// We are assuming that each file has one title as a H1 header...
	if raw[0] != '#' {
		return "", nil, fmt.Errorf("missing '#' (H1) at the start of file")
	}
	// ...followed by a line break and the contents
	for i, c := range raw {
		if c != '\n' {
			continue
		}

		title = string(bytes.TrimSpace(raw[1:i]))
		contents = bytes.TrimSpace(raw[i:])
		break
	}
	// If we scan the whole file and title or contents are empty, something
	// is wrong with the file
	if title == "" {
		return title, contents, fmt.Errorf("could not find title")
	}
	if contents == nil {
		return title, contents, fmt.Errorf("could not find content")
	}

	return title, contents, nil
}

func getAndValidateSlug(mdFilename, title string) (string, error) {
	// 01-my-awesome-blog-post.md => my-awesome-blog-post
	filenameSlug := strings.TrimSuffix(mdFilename[3:], ".md")
	// My awesome blog post => my-awesome-blog-post
	titleSlug := getSlug(title)

	// Add any filename that are known to be broken, generally because the
	// title got changed after publishing
	knownBrokenFilenames := []string{
		// Typo, should be "troubleshooting"
		"01-troubleshoting-zsh-lag-and-solutions-with-nix.md",
	}

	if filenameSlug != titleSlug && !slices.Contains(knownBrokenFilenames, mdFilename) {
		return filenameSlug, fmt.Errorf(
			"got conflicting slugs: filename slug: %s, title slug: %s",
			filenameSlug,
			titleSlug,
		)
	}

	return filenameSlug, nil
}

func getSlug(s string) string {
	return slug.Make(s)
}

func grabPosts(root string) (posts, error) {
	posts := orderedmap.NewOrderedMap[path, post]()

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
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

		// Get the base directory of the file
		dir := filepath.Base(filepath.Dir(path))
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
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf(
				"something went wrong when reading file: %s, error: %w",
				path,
				err,
			)
		}

		title, contents, err := extractTitleAndContents(raw)
		if err != nil || title == "" || contents == nil {
			return fmt.Errorf(
				"something is wrong with file: %s, error: %w",
				path,
				err,
			)
		}

		slug, err := getAndValidateSlug(d.Name(), title)
		if err != nil {
			return fmt.Errorf(
				"something is wrong with slug for file: %s, error: %w",
				path,
				err,
			)
		}

		if date.After(time.Now()) {
			log.Printf("[INFO]: ignoring future post: %s\n", path)
			return nil
		}

		posts.Set(path, post{
			title:    title,
			slug:     slug,
			contents: contents,
			date:     date,
		})

		return nil
	})

	return posts, err
}

func genRss(ps posts) string {
	feed := &feeds.Feed{
		Title:       "kokada's blog",
		Description: "# dd if=/dev/urandom of=/dev/brain0",
	}
	md := goldmark.New(goldmark.WithExtensions(
		NewLinkRewriter(blogMainUrl, nil),
		extension.GFM,
		highlighting.NewHighlighting(highlighting.WithStyle(highlightStyle)),
	))

	var items []*feeds.Item
	for path, post := range ps.ReverseIterator() {
		link := must1(url.JoinPath(blogMainUrl, path))
		var buf bytes.Buffer
		must(md.Convert(post.contents, &buf))
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

func genReadme(ps posts) string {
	var titles []string
	for path, post := range ps.ReverseIterator() {
		title := fmt.Sprintf(
			"- [%s](%s) - %s",
			post.title,
			path,
			post.date.Format(time.DateOnly),
		)
		titles = append(titles, title)
	}
	return fmt.Sprintf(readmeTemplate, strings.Join(titles, "\n"))
}

func main() {
	slugify := flag.String("slugify", "", "Slugify input (e.g.: for blog titles)")
	rss := flag.Bool("rss", false, "Generate RSS (XML) instead of README.md")
	prepare := flag.Bool("prepare", false, "Prepare and print posts to Mataroa (mostly for debug)")
	publish := flag.Bool("publish", false, "Publish updates to Maratoa instance")
	flag.Parse()

	if *slugify != "" {
		fmt.Println(getSlug(*slugify))
		os.Exit(0)
	}

	posts := must1(grabPosts("posts"))
	if *prepare {
		for filename, post := range prepareToMataroa(posts).Iterator() {
			fmt.Printf(
				"%s\n# %s\n\n%s\n",
				filename,
				post.title,
				post.contents,
			)
		}
	} else if *publish {
		publishToMataroa(posts)
	} else if *rss {
		fmt.Print(genRss(posts))
	} else {
		fmt.Print(genReadme(posts))
	}
}
