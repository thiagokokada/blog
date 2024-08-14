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

DAY := $(shell date)
_PARSED_DAY := $(shell date '+%Y-%m-%d' -d '$(DAY)')
.PHONY: day
day:
	mkdir -p '$(_PARSED_DAY)'

TITLE = $(error TITLE is not defined)
.PHONY: post
post: blog day
	DAY=$(_PARSED_DAY) ./.scripts/gen-post.sh

TITLE = $(error FILE is not defined)
.PHONY: draft
draft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '.$(notdir $(FILE))'

.PHONY: undraft
undraft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '$(patsubst .%,%,$(notdir $(FILE)))'

.PHONY: clean
clean:
	rm -rf blog
