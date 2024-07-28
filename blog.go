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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
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

const mataroaApiUrl = "https://capivaras.dev/api/"
const rssBaseUrl = "https://github.com/thiagokokada/blog/blob/main"
const readmeTemplate = `# Blog

Mirror of my blog in https://kokada.capivaras.dev/.

## Posts

[![RSS](https://img.shields.io/badge/RSS-FFA562?style=for-the-badge&logo=rss&logoColor=white)](https://raw.githubusercontent.com/thiagokokada/blog/main/rss.xml)

%s
`

var mataroaToken = os.Getenv("MATAROA_TOKEN")

type post struct {
	title    string
	file     string
	slug     string
	contents []byte
	date     time.Time
}

// https://capivaras.dev/api/docs/
type mataroaResponse struct {
	Ok    bool   `json:"ok"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Slug  string `json:"slug"`
	// Error string `json:"error"`
}

type mataroaPostRequest struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
}

type mataroaPatchRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Body        string `json:"body"`
	PublishedAt string `json:"published_at"`
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
			contents = raw[i:]
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
		// Find markdown files, but ignore hidden files
		if filepath.Ext(d.Name()) == ".md" && d.Name()[0] != '.' {
			dir := filepath.Dir(path)
			// Ignore root directory
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
		}
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

func mustGetMataroaPost(post post) (p mataroaResponse, r *http.Response) {
	url := must1(url.JoinPath(mataroaApiUrl, "posts", post.slug, "/"))
	req := must1(http.NewRequest("GET", url, nil))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mataroaToken))

	resp := must1(http.DefaultClient.Do(req))
	body := must1(io.ReadAll(resp.Body))
	json.Unmarshal(body, &p)
	return p, resp
}

func mustPatchMataroaPost(post post) (p mataroaResponse, r *http.Response) {
	url := must1(url.JoinPath(mataroaApiUrl, "posts", post.slug, "/"))
	reqBody := must1(json.Marshal(mataroaPatchRequest{
		Title:       post.title,
		Body:        string(post.contents),
		Slug:        post.slug,
		PublishedAt: post.date.Format(time.DateOnly),
	}))
	req := must1(http.NewRequest("PATCH", url, bytes.NewBuffer(reqBody)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mataroaToken))

	resp := must1(http.DefaultClient.Do(req))
	body := must1(io.ReadAll(resp.Body))
	json.Unmarshal(body, &p)
	return p, resp
}

func mustPostMataroaPost(post post) (p mataroaResponse, r *http.Response) {
	url := must1(url.JoinPath(mataroaApiUrl, "posts", "/"))
	reqBody := must1(json.Marshal(mataroaPostRequest{
		Title:       post.title,
		Body:        string(post.contents),
		PublishedAt: post.date.Format(time.DateOnly),
	}))
	req := must1(http.NewRequest("POST", url, bytes.NewBuffer(reqBody)))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mataroaToken))

	resp := must1(http.DefaultClient.Do(req))
	body := must1(io.ReadAll(resp.Body))
	json.Unmarshal(body, &p)
	return p, resp
}

func publishToMataroa(posts []post) {
	for _, post := range posts {
		p, resp := mustGetMataroaPost(post)
		if p.Ok {
			p, resp := mustPatchMataroaPost(post)
			fmt.Printf("[UPDATED] (code=%d): %+v\n", resp.StatusCode, p)
		} else if resp.StatusCode == 404 {
			p, resp := mustPostMataroaPost(post)
			fmt.Printf("[NEW] (code=%d): %+v\n", resp.StatusCode, p)
		} else {
			fmt.Printf("[ERROR] %s: %+v\n", post.slug, resp)
		}
	}
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
