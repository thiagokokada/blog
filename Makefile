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
.PHONY: day
day:
	mkdir -p '$(DATE)'

TITLE = $(error TITLE is not defined)
SLUG = $(shell ./blog -slugify "$(TITLE)")
.PHONY: post
post: blog day
	$(EDITOR) $(shell DATE=$(DATE) SLUG=$(SLUG) ./.scripts/gen-post.sh)

FILE = $(error FILE is not defined)
.PHONY: draft
draft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '.$(notdir $(FILE))'

.PHONY: undraft
undraft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '$(patsubst .%,%,$(notdir $(FILE)))'

.PHONY: image
image:
	@echo '[![$(DESCRIPTION)](/$(FILE))](/$(FILE))'

.PHONY: words
words:
	wc -w **/*.md

.PHONY: clean
clean:
	rm -rf blog
