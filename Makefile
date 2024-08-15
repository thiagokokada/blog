MARKDOWN := $(wildcard ./**/*.md)

.PHONY: all
all: README.md rss.xml

blog: *.go go.*
	go build

README.md: blog $(MARKDOWN)
	./blog > README.md

rss.xml: blog $(MARKDOWN)
	./blog -rss > rss.xml

.PHONY: publish
publish: blog
	./blog -publish

DAY := $(shell date '+%Y-%m-%d')
# -d is a GNUism, but sadly -j/-f is a BSDism
# If GNU date extensions doesn't work, just do not try to parse the DAY
_PARSED_DAY := $(shell date '+%Y-%m-%d' -d '$(DAY)' 2>/dev/null || echo '$(DAY)')
.PHONY: day
day:
	mkdir -p '$(_PARSED_DAY)'

TITLE = $(error TITLE is not defined)
.PHONY: post
post: blog day
	DAY=$(_PARSED_DAY) ./.scripts/gen-post.sh

FILE = $(error FILE is not defined)
.PHONY: draft
draft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '.$(notdir $(FILE))'

.PHONY: undraft
undraft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '$(patsubst .%,%,$(notdir $(FILE)))'

.PHONY: words
words:
	wc --words **/*.md

.PHONY: clean
clean:
	rm -rf blog
