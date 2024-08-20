module github.com/thiagokokada/blog

go 1.23

require (
	github.com/elliotchance/orderedmap/v2 v2.4.0
	github.com/gorilla/feeds v1.2.0
	github.com/gosimple/slug v1.14.0
	github.com/teekennedy/goldmark-markdown v0.3.0
	github.com/yuin/goldmark v1.7.4
	github.com/yuin/goldmark-highlighting v0.0.0-20220208100518-594be1970594
)

require (
	github.com/alecthomas/chroma v0.10.0 // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
)

replace github.com/teekennedy/goldmark-markdown => github.com/thiagokokada/goldmark-markdown v0.0.0-20240820111219-f30775d8ed15
