MARKDOWN := $(shell find . -type f -name '*.md' -not -name README.md)

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
	mkdir -p $(shell date '+%Y-%m-%d')

.PHONE: draft
draft:
	mv "$(FILE)" "$(dir $(FILE)).$(notdir $(FILE))"

.PHONE: undraft
undraft:
	mv "$(FILE)" "$(dir $(FILE))$(patsubst .%,%,$(notdir $(FILE)))"

.PHONY: clean
clean:
	rm -rf blog
