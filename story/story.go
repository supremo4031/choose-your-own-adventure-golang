package story

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var defaultHandlerTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Choose Your Own Adventure</title>
</head>
<body>
  <h1>{{.Title}}</h1>
  {{range .Paragraphs}}
    <p>{{.}}</p>
  {{end}}
  <ul>
    {{range .Options}}
      <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
    {{end}}
  </ul>
</body>
</html>`

type Story map[string]Chapter

type Chapter struct {
	Title string `json:"title"`
	Paragraphs []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Chapter string `json:"arc"`
}

func JsonStory(r io.Reader) (Story, error) {
	decoder := json.NewDecoder(r)
	var story Story
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

type handler struct {
	s Story
}

func NewHandler(s Story) handler {
	return handler{s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Parse(defaultHandlerTemplate))
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	// "/intro" => "intro"
	path = path[1:]

	// ["intro"]
	if chapter, ok := h.s[path]; ok {
		err := tmpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
	}
}

