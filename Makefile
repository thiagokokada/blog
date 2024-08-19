.PHONY: all
all: README.md rss.xml

blog: *.go go.*
	go build -v

MARKDOWN := $(wildcard ./**/*.md)
README.md: blog $(MARKDOWN)
	./blog > README.md

rss.xml: blog $(MARKDOWN)
	./blog -rss > rss.xml

.PHONY: publish
publish: blog
	./blog -publish

DATE := $(shell date '+%Y-%m-%d')
# -d is a GNUism, but sadly -j/-f is a BSDism
# If GNU date extensions doesn't work, just do not try to parse the DATE
_PARSED_DATE := $(shell date '+%Y-%m-%d' -d '$(DATE)' 2>/dev/null || echo '$(DATE)')
.PHONY: day
day:
	mkdir -p '$(_PARSED_DATE)'

TITLE = $(error TITLE is not defined)
.PHONY: post
post: blog day
	DATE=$(_PARSED_DATE) ./.scripts/gen-post.sh

FILE = $(error FILE is not defined)
.PHONY: draft
draft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '.$(notdir $(FILE))'

.PHONY: undraft
undraft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '$(patsubst .%,%,$(notdir $(FILE)))'

.PHONY: words
words:
	wc -w **/*.md

.PHONY: clean
clean:
	rm -rf blog
