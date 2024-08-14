MARKDOWN := $(shell find . -type f -name '*.md' -not -name README.md)

.PHONY: all
all: README.md rss.xml

blog: *.go go.*
	go build

README.md: blog $(MARKDOWN)
	./blog > README.md

rss.xml: blog $(MARKDOWN)
	./blog -rss > rss.xml 2>/dev/null

.PHONY: publish
publish: blog
	./blog -publish

DAY := $(shell date)
_PARSED_DAY := $(shell date '+%Y-%m-%d' -d '$(DAY)')
.PHONY: day
day:
	mkdir -p '$(_PARSED_DAY)'

.PHONY: post
post: blog day
	@[ "${TITLE}" ] || ( echo ">> TITLE is not set"; exit 1 )
	./.scripts/gen-post.sh '$(_PARSED_DAY)' '$(TITLE)'

.PHONY: draft
draft:
	@[ "${FILE}" ] || ( echo ">> FILE is not set"; exit 1 )
	mv '$(FILE)' '$(dir $(FILE)).$(notdir $(FILE))'

.PHONY: undraft
undraft:
	@[ "${FILE}" ] || ( echo ">> FILE is not set"; exit 1 )
	mv '$(FILE)' '$(dir $(FILE))$(patsubst .%,%,$(notdir $(FILE)))'

.PHONY: clean
clean:
	rm -rf blog
