package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Payload struct {
	Query      string        `json:"query"`
	Tags       []interface{} `json:"tags"`
	Sort       []string      `json:"sort"`
	Types      []interface{} `json:"types"`
	Categories []interface{} `json:"categories"`
	StartDate  time.Time     `json:"startDate"`
	EndDate    time.Time     `json:"endDate"`
	From       int           `json:"from"`
}

type Hits struct {
	Score       interface{} `json:"score"`
	ID          int         `json:"id"`
	UUID        string      `json:"uuid"`
	Views       int         `json:"views"`
	Name        string      `json:"name"`
	FullSlug    string      `json:"full_slug"`
	PublishedAt time.Time   `json:"published_at"`
}

type Response struct {
	Total     int      `json:"total"`
	Uuids     []string `json:"uuids"`
	Storyblok string   `json:"storyblok"`
	Principal struct {
		Score       interface{} `json:"score"`
		ID          int         `json:"id"`
		UUID        string      `json:"uuid"`
		Views       int         `json:"views"`
		Name        string      `json:"name"`
		FullSlug    string      `json:"full_slug"`
		PublishedAt time.Time   `json:"published_at"`
	} `json:"principal"`
	Hits []Hits `json:"hits"`
}

type Result struct {
	Data []Hits `json:"data"`
}

func fetchLupaUrl(query string, currentPage int) (Response, error) {
	var responseApi = Response{}

	layout := time.RFC3339[:len("yyyy-MM-dd'T'HH:mm:ssZ")]
	startDate, err := time.Parse(layout, "2015-11-01T00:00:01.000Z")
	if err != nil {
		return responseApi, err
	}

	timeNow := time.Now()
	endDate, err := time.Parse(layout, timeNow.Format(layout))
	if err != nil {
		return responseApi, err
	}

	data := Payload{
		Query: query,
		Tags:  nil,
		Sort: []string{
			"published_at",
			"desc",
		},
		Types:      nil,
		Categories: nil,
		StartDate:  startDate,
		EndDate:    endDate,
		From:       currentPage * 10, // ele tras 10 resultados por vez
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return responseApi, err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.lupa.news/v1/search", body)
	if err != nil {
		return responseApi, err
	}

	req.Header.Set("Authority", "api.lupa.news")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://lupa.uol.com.br")
	req.Header.Set("Referer", "https://lupa.uol.com.br/")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return responseApi, err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	json.Unmarshal(bodyBytes, &responseApi)

	return responseApi, nil
}

type byDateDesc []Hits

func (s byDateDesc) Len() int {
	return len(s)
}
func (s byDateDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byDateDesc) Less(i, j int) bool {
	return s[i].PublishedAt.After(s[j].PublishedAt)
}

func CrawlLupa(query string) (Result, error) {
	responseApi, err := fetchLupaUrl(query, 0)
	if err != nil {
		return Result{}, err
	}

	if responseApi.Total == 0 {
		return Result{}, err
	}

	pages := responseApi.Total / len(responseApi.Hits)

	var results Result

	results.Data = append(results.Data, responseApi.Hits...)

	var wg sync.WaitGroup

	for i := 1; i < pages; i++ {
		wg.Add(1)

		i := i
		go func() {
			defer wg.Done()
			result, _ := fetchLupaUrl(query, i)
			results.Data = append(results.Data, result.Hits...)
		}()
	}

	wg.Wait()

	sort.Sort(byDateDesc(results.Data))
	return results, nil
}
