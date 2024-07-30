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
)

const mataroaApiUrl = "https://capivaras.dev/api/"

var mataroaToken = os.Getenv("MATAROA_TOKEN")

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

func mustMataroaReq(method string, elem []string, body []byte) (p mataroaResponse, r *http.Response) {
	// generate a Mataroa URL, ensure '/' at the end
	reqUrl := must1(url.JoinPath(mataroaApiUrl, elem...))
	reqUrl = must1(url.JoinPath(reqUrl, "/"))

	// Prepare request payload if non-nil
	var reqBuf io.Reader
	if body != nil {
		reqBuf = bytes.NewBuffer(body)
	}
	req := must1(http.NewRequest(method, reqUrl, reqBuf))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", mataroaToken))

	// Do request and return response
	resp := must1(http.DefaultClient.Do(req))
	respBody := must1(io.ReadAll(resp.Body))
	json.Unmarshal(respBody, &p)
	return p, resp
}

func mustGetMataroaPost(post post) (p mataroaResponse, r *http.Response) {
	return mustMataroaReq("GET", []string{"posts", post.slug}, nil)
}

func mustPatchMataroaPost(post post) (p mataroaResponse, r *http.Response) {
	reqBody := must1(json.Marshal(mataroaPatchRequest{
		Title:       post.title,
		Body:        string(post.contents),
		Slug:        post.slug,
		PublishedAt: post.date.Format(time.DateOnly),
	}))
	return mustMataroaReq("PATCH", []string{"posts", post.slug}, reqBody)
}

func mustPostMataroaPost(post post) (p mataroaResponse, r *http.Response) {
	reqBody := must1(json.Marshal(mataroaPostRequest{
		Title:       post.title,
		Body:        string(post.contents),
		PublishedAt: post.date.Format(time.DateOnly),
	}))
	return mustMataroaReq("POST", []string{"posts"}, reqBody)
}

func publishToMataroa(posts []post) {
	for _, post := range posts {
		p, resp := mustGetMataroaPost(post)
		if p.Ok {
			p, resp = mustPatchMataroaPost(post)
			log.Printf("[UPDATED] (code=%d): %+v\n", resp.StatusCode, p)
		} else if resp.StatusCode == 404 {
			p, resp = mustPostMataroaPost(post)
			log.Printf("[NEW] (code=%d): %+v\n", resp.StatusCode, p)
		} else {
			log.Printf("[ERROR] %s: %+v\n", post.slug, resp)
		}

		if resp.StatusCode != 200 {
			panic(fmt.Sprintf("non-200 response for post: %s", post.slug))
		}
	}
}
