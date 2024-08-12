MARKDOWN := $(shell find . -type f -name '*.md' -not -name README.md)
TODAY := $(shell date '+%Y-%m-%d')

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

.PHONY: today
today:
	mkdir -p $(TODAY)


.PHONY: post
post: blog today
	@[ "${TITLE}" ] || ( echo ">> TITLE is not set"; exit 1 )
	./.scripts/gen-post.sh $(TODAY) "$(TITLE)"

.PHONE: draft
draft:
	@[ "${FILE}" ] || ( echo ">> FILE is not set"; exit 1 )
	mv "$(FILE)" "$(dir $(FILE)).$(notdir $(FILE))"

.PHONE: undraft
undraft:
	@[ "${FILE}" ] || ( echo ">> FILE is not set"; exit 1 )
	mv "$(FILE)" "$(dir $(FILE))$(patsubst .%,%,$(notdir $(FILE)))"

.PHONY: clean
clean:
	rm -rf blog
