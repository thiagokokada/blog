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
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// linkRewriter is the main struct for your extension
type linkRewriter struct {
	prefixUrl string
	posts     posts
}

// NewLinkRewriter returns a new instance of LinkRewriter
func NewLinkRewriter(prefixUrl string, posts posts) *linkRewriter {
	return &linkRewriter{prefixUrl, posts}
}

// Extend will be called by Goldmark to add your extension
func (e *linkRewriter) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(util.Prioritized(e, 0)))
}

// Transform is the method that modifies the AST
func (e *linkRewriter) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if link, ok := n.(*ast.Link); ok {
			e.rewriteLink(link)
		}
		if image, ok := n.(*ast.Image); ok {
			e.rewriteImage(image)
		}
		return ast.WalkContinue, nil
	})
}

func hasAnyExtension(s string, extensions ...string) bool {
	ext := filepath.Ext(s)
	if ext == "" {
		return false
	}
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}

// rewriteLink modifies the link URL
func (e *linkRewriter) rewriteLink(l *ast.Link) {
	link := string(l.Destination)

	if strings.HasPrefix(link, ".") {
		log.Printf("[WARN]: relative link reference found: %s\n", link)
	}

	if strings.HasPrefix(link, "/") {
		var dest string

		if hasAnyExtension(link, ".png", ".jpg", ".jpeg") {
			// If the link is an image, we will point it to
			// blogRawUrl
			if _, err := os.Stat(filepath.Join(".", link)); err == nil {
				dest = must1(url.JoinPath(blogRawUrl, link))
			} else {
				log.Printf("[WARN] did not find image: %s\n", link)
				return
			}
		} else if e.posts != nil {
			// If posts are not nil, it means we will grab the slug
			// from posts
			post, ok := e.posts.Get(link[1:])
			if ok {
				dest = must1(url.JoinPath(e.prefixUrl, post.slug))
			} else {
				log.Printf("[WARN] did not find reference to link: %s\n", link)
				return
			}
		} else {
			// Else we will just append the prefixUrl to the link
			if _, err := os.Stat(filepath.Join(".", link)); err == nil {
				dest = must1(url.JoinPath(e.prefixUrl, link))
			} else {
				log.Printf("[WARN] did not find link: %s\n", link)
				return
			}
		}

		l.Destination = []byte(dest)
	}
}

// rewriteImage modifies the image URL
func (e *linkRewriter) rewriteImage(i *ast.Image) {
	image := string(i.Destination)

	if strings.HasPrefix(image, ".") {
		log.Printf("[WARN]: relative image link reference found: %s\n", image)
	}

	if strings.HasPrefix(image, "/") {
		dest := must1(url.JoinPath(blogRawUrl, image))
		i.Destination = []byte(dest)
	}
}
