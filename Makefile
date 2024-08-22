export POST_ROOT := posts
export DATE := $(shell date '+%Y-%m-%d')

MARKDOWN := $(wildcard $(POST_ROOT)/**/*.md)
TITLE = $(error TITLE is not defined)
FILE = $(error FILE is not defined)
SLUG = $(shell ./blog -slugify "$(TITLE)")

.PHONY: all
all: README.md rss.xml

blog: *.go go.* vendor
	go build -v

README.md: blog $(MARKDOWN)
	./blog > README.md

rss.xml: blog $(MARKDOWN)
	./blog -rss > rss.xml

.PHONY: publish
publish: blog
	./blog -publish

.PHONY: day
day:
	mkdir -p '$(POST_ROOT)/$(DATE)'

.PHONY: post
post: blog day
	$(EDITOR) $(shell SLUG=$(SLUG) ./.scripts/gen-post.sh)

.PHONY: draft
draft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '.$(notdir $(FILE))'

.PHONY: undraft
undraft:
	cd '$(dir $(FILE))' && mv '$(notdir $(FILE))' '$(patsubst .%,%,$(notdir $(FILE)))'

.PHONY: image
image:
	@echo '[![$(DESCRIPTION)](/$(POST_ROOT)/$(FILE))](/$(FILE))'

.PHONY: words
words:
	wc -w $(POST_ROOT)/**/*.md

.PHONY: clean
clean:
	rm -rf blog
