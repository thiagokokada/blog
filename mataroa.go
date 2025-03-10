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

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/elliotchance/orderedmap/v3"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

const (
	mataroaBaseUrl = "https://capivaras.dev"
	mataroaApiUrl  = mataroaBaseUrl + "/api/"
	mataroaBlogUrl = "https://kokada.dev/blog/"
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

func mataroaReq(method string, url string, body []byte) (m mataroaResponse, r *http.Response, err error) {
	// Prepare request payload if non-nil
	var reqBuf io.Reader
	if body != nil {
		reqBuf = bytes.NewBuffer(body)
	}
	req := must1(http.NewRequest(method, url, reqBuf))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mataroaToken))

	// Do request and return response
	r, err = http.DefaultClient.Do(req)
	if err != nil {
		return m, r, fmt.Errorf("Mataroa request error: %w", err)
	}

	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		return m, r, fmt.Errorf("Mataroa response error: %w", err)
	}

	err = json.Unmarshal(rBody, &m)
	if err != nil {
		return m, r, fmt.Errorf("Mataroa JSON unmarshal error: %w", err)
	}

	return m, r, nil
}

func getMataroaPost(slug string) (mataroaResponse, *http.Response, error) {
	return mataroaReq("GET", mustMataroaUrl("posts", slug), nil)
}

func patchMataroaPost(slug string, p post) (mataroaResponse, *http.Response, error) {
	reqBody := must1(json.Marshal(mataroaPatchRequest{
		Title:       p.title,
		Body:        string(p.contents),
		Slug:        p.slug,
		PublishedAt: p.date.Format(time.DateOnly),
	}))
	return mataroaReq("PATCH", mustMataroaUrl("posts", slug), reqBody)
}

func postMataroaPost(p post) (mataroaResponse, *http.Response, error) {
	reqBody := must1(json.Marshal(mataroaPostRequest{
		Title:       p.title,
		Body:        string(p.contents),
		PublishedAt: p.date.Format(time.DateOnly),
	}))
	return mataroaReq("POST", mustMataroaUrl("posts"), reqBody)
}

func prepareToMataroa(ps posts) posts {
	md := goldmark.New(
		goldmark.WithExtensions(
			NewLinkRewriter(mataroaBlogUrl, ps),
			extension.GFM,
			highlighting.NewHighlighting(
				// No style since we are reusing the style from
				// Mataroa
				highlighting.WithFormatOptions(html.WithClasses(true)),
			),
		),
	)

	preparedPosts := orderedmap.NewOrderedMap[path, post]()
	for el := ps.Front(); el != nil; el = el.Next() {
		path, post := el.Key, el.Value
		buf := bytes.Buffer{}
		must(md.Convert([]byte(post.contents), &buf))
		post.contents = buf.Bytes()
		preparedPosts.Set(path, post)
	}
	return preparedPosts
}

func publishToMataroa(ps posts) {
	if mataroaToken == "" {
		log.Fatal("empty MATAROA_TOKEN environment variable")
	}

	for post := range prepareToMataroa(ps).Values() {
		p, resp := must2(getMataroaPost(post.slug))
		var err error
		if p.Ok {
			p, resp, err = patchMataroaPost(post.slug, post)
			log.Printf("[UPDATED] (code=%d): %+v\n", resp.StatusCode, p)
		} else if resp.StatusCode == 404 {
			p, resp, err = postMataroaPost(post)
			log.Printf("[NEW] (code=%d): %+v\n", resp.StatusCode, p)

			if p.Ok && p.Slug != post.slug {
				log.Printf(
					"[INFO] Updating slug since they're different, Mataroa slug: %s, generated one: %s",
					p.Slug,
					post.slug,
				)
				p, resp, err = patchMataroaPost(p.Slug, post)
				log.Printf("[UPDATED] (code=%d): %+v\n", resp.StatusCode, p)
			}
		}

		if resp.StatusCode != 200 {
			err = fmt.Errorf(
				"non-200 (code=%d) status code for post=%s, response: %+v",
				resp.StatusCode,
				post.slug,
				resp,
			)
		}

		must(err)
	}
}
