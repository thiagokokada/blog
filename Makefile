MARKDOWN := $(shell find . -type f -name '*.md' -not -name README.md)

.PHONY: all
all: README.md rss.xml

blog: blog.go go.*
	go build

README.md: blog $(MARKDOWN)
	./blog > README.md

rss.xml: blog $(MARKDOWN)
	./blog -rss > rss.xml

.PHONY: publish
publish: blog
	./blog -publish

.PHONY: clean
clean:
	rm -rf blog
