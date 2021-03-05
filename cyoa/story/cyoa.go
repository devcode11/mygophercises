package story

import (
	"encoding/json"
	"io"
	"net/http"
	"html/template"
	"log"
	"strings"
)

var htmlStoryTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<title>CYOA</title>
	</head>
<body>
	<h1>{{ .Title }}</h1>
	{{ range .Stories }}
	<p>{{ . }}</p>
	{{end}}
	<ul>
	{{ range .Options }}
	<li><a href="/{{ .Arc }}">{{.Text}}</a></li>
	{{end}}
	</ul>
</body>
</html>`

type StoryArc struct {
	Title   string   `json:"title"`
	Stories []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Story map[string]StoryArc

func FromJson(r io.Reader) (Story, error) {
	var story Story
	d := json.NewDecoder(r)
	err := d.Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}

type StoryHandler struct {
	Story Story
}

func (h StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storyData, ok := h.Story[strings.TrimPrefix(r.URL.Path, "/")]
	if !ok {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.New("storyArcTmpl").Parse(htmlStoryTmpl))
	// log.Println(r.URL.Path, sd)
	err := tmpl.Execute(w, storyData)
	if err != nil {
		log.Printf("Path %s : %v",r.URL.Path, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}