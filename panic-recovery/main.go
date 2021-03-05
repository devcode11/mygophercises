package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	// "html"

	"main/abc"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

const srcFilePrefix = "/src/"
const lineLinkPrefix = "line-"

var srcFileLine = regexp.MustCompile(`[\t]([^:]+):(\d+)`)

const (
	serverEnv = "SERVERENV"
	devEnv = "DEV"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexPage)
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/abc/", abcHello)
	mux.HandleFunc("/hello", hello)
	if os.Getenv(serverEnv) == devEnv {
		mux.HandleFunc(srcFilePrefix, srcFileHander)
	}
	log.Fatal(http.ListenAndServe(":3000", recoverMiddleware(mux)))
}

func srcFileHander(w http.ResponseWriter, r *http.Request) {
	splitRes := strings.SplitAfterN(r.URL.Path, srcFilePrefix, 2)
	filePath := splitRes[1]
	if !strings.HasSuffix(filePath, ".go") {
		log.Println("Cannot render", filePath)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Cannot render", filePath)
		return
	}
	var line int
	v := r.FormValue("line")
	if l, err := strconv.Atoi(v); err == nil {
		line = l
	}

	contents, err := ioutil.ReadFile("./" + filePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Could not find", filePath)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=UTF-8")

	// fmt.Fprintln(w, "<code><pre>")
	// fmt.Fprintln(w, html.EscapeString(string(contents)))
	// fmt.Fprintln(w, "</pre></code>")

	// err = quick.Highlight(w, string(contents), "go", "html", "monokai")
	formatter := html.New(html.WithLineNumbers(true), html.HighlightLines([][2]int{{line, line}}), html.LineNumbersInTable(true), html.LinkableLineNumbers(true, "line-"))
	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	iterator, err := lexer.Tokenise(nil, string(contents))
	err = formatter.Format(w, style, iterator)
	if err != nil {
		log.Println(err)
	}
}

func formatStackTrace(trace string) string {

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	matches := srcFileLine.FindAllStringSubmatch(trace, -1)
	for _, match := range matches {
		fl := match[1]
		if !strings.HasPrefix(fl, pwd) {
			continue
		}
		flln := match[1] + ":" + match[2]
		link := strings.Replace(fl, pwd+"/", srcFilePrefix, 1)
		link = strings.ReplaceAll(link, " ", "%20")
		trace = strings.Replace(trace, flln, fmt.Sprintf("<a href=%s?line=%s#%s%s>%s</a>", link, match[2], lineLinkPrefix, match[2], flln), -1)
	}
	return trace
}

func recoverMiddleware(app http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Add("Content-Type", "text/html; charset=UTF-8")
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h2>%v</h2>", err)
				if os.Getenv(serverEnv) == devEnv {
					stackTrace := string(debug.Stack())
					linkedStackTrace := formatStackTrace(stackTrace)
					fmt.Fprintf(w, "<pre>%s</pre>", linkedStackTrace)
				}
			}
		}()

		app.ServeHTTP(w, r)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}

func abcHello(w http.ResponseWriter, r *http.Request) {
	abc.Hello()
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintln(w, "Visit <a href=\"/panic\">panic</a>")
}