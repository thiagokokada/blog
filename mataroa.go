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

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

const (
	mataroaBaseUrl = "https://kokada.capivaras.dev"
	mataroaApiUrl  = mataroaBaseUrl + "/api/"
	mataroaBlogUrl = mataroaBaseUrl + "/blog/"
)

var (
	mataroaToken = os.Getenv("MATAROA_TOKEN")
	md           goldmark.Markdown
)

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

func mustMataroaReq(method string, url string, body []byte) (p mataroaResponse, r *http.Response) {
	// Prepare request payload if non-nil
	var reqBuf io.Reader
	if body != nil {
		reqBuf = bytes.NewBuffer(body)
	}
	req := must1(http.NewRequest(method, url, reqBuf))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mataroaToken))

	// Do request and return response
	resp := must1(http.DefaultClient.Do(req))
	respBody := must1(io.ReadAll(resp.Body))
	json.Unmarshal(respBody, &p)
	return p, resp
}

func mustGetMataroaPost(slug string) (p mataroaResponse, r *http.Response) {
	return mustMataroaReq("GET", mustMataroaUrl("posts", slug), nil)
}

func mustPatchMataroaPost(slug string, post post) (p mataroaResponse, r *http.Response) {
	buf := bytes.Buffer{}
	must(md.Convert([]byte(post.contents), &buf))
	reqBody := must1(json.Marshal(mataroaPatchRequest{
		Title:       post.title,
		Body:        buf.String(),
		Slug:        post.slug,
		PublishedAt: post.date.Format(time.DateOnly),
	}))
	return mustMataroaReq("PATCH", mustMataroaUrl("posts", slug), reqBody)
}

func mustPostMataroaPost(post post) (p mataroaResponse, r *http.Response) {
	buf := bytes.Buffer{}
	must(md.Convert([]byte(post.contents), &buf))
	reqBody := must1(json.Marshal(mataroaPostRequest{
		Title:       post.title,
		Body:        buf.String(),
		PublishedAt: post.date.Format(time.DateOnly),
	}))
	return mustMataroaReq("POST", mustMataroaUrl("posts"), reqBody)
}

func publishToMataroa(posts posts) {
	if mataroaToken == "" {
		log.Fatal("empty MATAROA_TOKEN environment variable")
	}

	md = goldmark.New(
		goldmark.WithRenderer(markdown.NewRenderer()),
		goldmark.WithExtensions(
			NewLinkRewriter(mataroaBlogUrl, posts),
			extension.GFM,
		),
	)
	for _, post := range posts.Iterator() {
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
