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
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

const (
	mataroaBaseUrl = "https://kokada.capivaras.dev"
	mataroaApiUrl  = mataroaBaseUrl + "/api/"
	mataroaBlogUrl = mataroaBaseUrl + "/blog/"
)

var mataroaToken = os.Getenv("MATAROA_TOKEN")

// https://capivaras.dev/api/docs/
type mataroaResponse struct {
	Ok    bool   `json:"ok"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Slug  string `json:"slug"`
	Error string `json:"error"`
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

func mustMataroaUrl(elem ...string) string {
	// generate a Mataroa URL, ensure '/' at the end
	mUrl := must1(url.JoinPath(mataroaApiUrl, elem...))
	mUrl = must1(url.JoinPath(mUrl, "/"))
	return mUrl
}

func mustMataroaReq(method string, url string, body []byte) (m mataroaResponse, r *http.Response) {
	// Prepare request payload if non-nil
	var reqBuf io.Reader
	if body != nil {
		reqBuf = bytes.NewBuffer(body)
	}
	req := must1(http.NewRequest(method, url, reqBuf))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mataroaToken))

	// Do request and return response
	r = must1(http.DefaultClient.Do(req))
	rBody := must1(io.ReadAll(r.Body))
	json.Unmarshal(rBody, &m)
	return m, r
}

func mustGetMataroaPost(slug string) (mataroaResponse, *http.Response) {
	return mustMataroaReq("GET", mustMataroaUrl("posts", slug), nil)
}

func mustPatchMataroaPost(slug string, p post) (mataroaResponse, *http.Response) {
	reqBody := must1(json.Marshal(mataroaPatchRequest{
		Title:       p.title,
		Body:        string(p.contents),
		Slug:        p.slug,
		PublishedAt: p.date.Format(time.DateOnly),
	}))
	return mustMataroaReq("PATCH", mustMataroaUrl("posts", slug), reqBody)
}

func mustPostMataroaPost(p post) (mataroaResponse, *http.Response) {
	reqBody := must1(json.Marshal(mataroaPostRequest{
		Title:       p.title,
		Body:        string(p.contents),
		PublishedAt: p.date.Format(time.DateOnly),
	}))
	return mustMataroaReq("POST", mustMataroaUrl("posts"), reqBody)
}

func prepareToMataroa(ps posts) posts {
	md := goldmark.New(
		goldmark.WithRenderer(
			markdown.NewRenderer(markdown.WithSubListLength(2)),
		),
		goldmark.WithExtensions(
			NewLinkRewriter(mataroaBlogUrl, ps),
			extension.GFM,
		),
	)

	preparedPosts := orderedmap.NewOrderedMap[path, post]()
	for filename, p := range ps.Iterator() {
		buf := bytes.Buffer{}
		must(md.Convert([]byte(p.contents), &buf))
		p.contents = buf.Bytes()
		preparedPosts.Set(filename, p)
	}
	return preparedPosts
}

func publishToMataroa(ps posts) {
	if mataroaToken == "" {
		log.Fatal("empty MATAROA_TOKEN environment variable")
	}

	for _, post := range prepareToMataroa(ps).Iterator() {
		p, resp := mustGetMataroaPost(post.slug)
		if p.Ok {
			p, resp = mustPatchMataroaPost(post.slug, post)
			log.Printf("[UPDATED] (code=%d): %+v\n", resp.StatusCode, p)
		} else if resp.StatusCode == 404 {
			p, resp = mustPostMataroaPost(post)
			log.Printf("[NEW] (code=%d): %+v\n", resp.StatusCode, p)

			if p.Ok && p.Slug != post.slug {
				log.Printf(
					"[INFO] Updating slug since they're different, Mataroa slug: %s, generated one: %s",
					p.Slug,
					post.slug,
				)
				p, resp = mustPatchMataroaPost(p.Slug, post)
				log.Printf("[UPDATED] (code=%d): %+v\n", resp.StatusCode, p)
			}
		}

		if resp.StatusCode != 200 {
			err := fmt.Errorf(
				"non-200 (code=%d) status code for post=%s, response: %+v",
				resp.StatusCode,
				post.slug,
				resp,
			)
			panic(err)
		}
	}
}
