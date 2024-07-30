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
		return ast.WalkContinue, nil
	})
}

// rewriteLink modifies the link URL
func (e *linkRewriter) rewriteLink(l *ast.Link) {
	link := string(l.Destination)

	if strings.HasPrefix(link, ".") {
		log.Printf("[WARN]: relative link reference found: %s\n", link)
	}

	if strings.HasPrefix(link, "/") {
		if e.posts != nil {
			// If posts are not nil, it means we will grab the slug
			// from posts
			post, ok := e.posts.Get(link[1:])
			if ok {
				dest := must1(url.JoinPath(e.prefixUrl, post.slug))
				l.Destination = []byte(dest)
			} else {
				log.Printf("[WARN]: did not find reference to link: %s\n", link)
			}
		} else {
			// Else we will just append the prefixUrl to the link
			dest := must1(url.JoinPath(e.prefixUrl, link))
			l.Destination = []byte(dest)
		}
	}
}
