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
	mkdir $(shell date '+%Y-%m-%d')

.PHONE: draft
draft:
	mv $(FILE) $(shell dirname $(FILE))/.$(shell basename $(FILE))

.PHONE: undraft
undraft:
	mv $(FILE) $(shell dirname $(FILE))/$(shell echo $(shell basename $(FILE)) | tail -c +2)

.PHONY: clean
clean:
	rm -rf blog
