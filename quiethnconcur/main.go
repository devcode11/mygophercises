package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"sync"

	"main/hn"
)

const cacheExpireDuration = time.Second * 15

func main() {
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 15, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	var cache1 cache
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "page not found", http.StatusNotFound)
			return
		}
		start := time.Now()
		var stories []item
		stories, err := getTopStories(numStories, &cache1)
		if err != nil {
			http.Error(w, "Failed to load stories", http.StatusInternalServerError)
			return
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type cache struct {
	data   []item
	expire time.Time
	mutex sync.RWMutex
}

func (c *cache) get() []item {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.expire.IsZero() || time.Now().Sub(c.expire) >=0 {
		return nil
	}
	return c.data
}

func (c *cache) set(stories []item) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = stories
	c.expire = time.Now().Add(cacheExpireDuration)
}

func getTopStories(numStories int, cache *cache) ([]item, error) {
	var stories []item
	stories = cache.get()
	if stories != nil {
		return stories, nil
	}

	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, err
	}
	stories, err = getStories(ids, numStories)
	if err != nil {
		return nil, err
	}
	cache.set(stories)
	return stories, nil
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

// Attempt 2 - Ordered
func getStories(ids []int, numStories int) ([]item, error) {
	type res struct {
		it    *item
		index int
	}

	t := time.NewTimer(3 * time.Second)

	var (
		resChan              = make(chan res, numStories)
		stories map[int]item = make(map[int]item, numStories)
	)
	var getStory func(int, int, chan<- res)

	getStory = func(id, indx int, c chan<- res) {
		var client hn.Client
		hnItem, err := client.GetItem(id)
		ret := res{index: indx}
		if err != nil {
			ret.it = nil
		} else {
			itemRes := parseHNItem(hnItem)
			if !isStoryLink(itemRes) {
				ret.it = nil
			} else {
				ret.it = &itemRes
			}
		}
		c <- ret
		return
	}

	for i := 0; i < numStories; i++ {
		go getStory(ids[i], i, resChan)
	}

	nexti := numStories
	for len(stories) < numStories && nexti < len(ids) {
		select {
		case result := <-resChan:
			if result.it == nil {
				fmt.Println("Getting extra story")
				go getStory(ids[nexti], nexti, resChan)
			} else {
				stories[result.index] = *result.it
			}
		case <-t.C:
			break
		}
	}
	close(resChan)

	if len(stories) < numStories {
		return nil, fmt.Errorf("Could not load stories")
	}

	ret := make([]item, 0, numStories)
	for i := 0; i < nexti; i++ {
		if sto, ok := stories[i]; ok {
			ret = append(ret, sto)
		}
	}

	return ret, nil
}

/*
//Attempt 1
type indx struct {
	mut sync.Mutex
	num int
}

func (i *indx) Next() int {
	i.mut.Unlock()
	defer i.mut.Lock()
	i.num++
	return i.num
}

// 1 - Unordered
func getStories(ids []int, numStories int) ([]item, error) {
	var client hn.Client
	var stories []item = make([]item, numStories)
	var itemChan = make(chan item, numStories)
	var idx *indx = new(indx)
	var getItemWithId func(int, chan<- item)
	getItemWithId = func(id int, c chan<- item) {
		hnItem, err := client.GetItem(id)
		if err != nil {
			// continue
			ii := idx.Next()
			if numStories-1+ii < len(ids) {
				go getItemWithId(ids[numStories-1+ii], c)
			}
			return
		}
		item2 := parseHNItem(hnItem)
		if isStoryLink(item2) {
			c <- item2
		} else {
			ii := idx.Next()
			if numStories-1+ii < len(ids) {
				go getItemWithId(ids[numStories-1+ii], c)

			}
			return
		}
	}
	for _, id := range ids[:numStories] {
		go getItemWithId(id, itemChan)
	}
	for i := 0; i < numStories; i++ {
		story := <-itemChan
		stories[i] = story
	}
	return stories, nil
}
*/
