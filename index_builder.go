package indexer

import (
  "bufio"
  "io/ioutil"
  "fmt"
  "os"
  "net/http"
  // "time"
  "log"
  "strings"
)

type httpResponse struct {
  title string
  text string
  err error
}


func makeRequest(url string, buffered_chan chan struct{}, results_chan chan <- *httpResponse) {
  buffered_chan <- struct{}{}

  resp, err := http.Get(url)
  // fmt.Println(err)
  body, _ := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()

  result := &httpResponse{"", string(body), err}
  splitted := strings.SplitN(result.text, "\n", 2)
  result.title = splitted[0]

  results_chan <- result

  <- buffered_chan
}

func runIndexer(filename string, filename2 string) {
  // Try to open text file with urls
  file, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  urls := readFile(file)

  res := loadReviews(urls, 100)

  index , counts_index := indexReviews(res, filename2)
  for _, val:= range counts_index["study"] {
    fmt.Printf(" val: %d\n", val)
  }

  res = nil
  return index, counts_index
}
func indexReviews(reviews []httpResponse, filter_filename string) (map[string] map[string] []int, map[string] map[string] int) {
  index := make(map[string] (map[string] []int))
  count_index := make(map[string] (map[string] int))
  replacer := strings.NewReplacer(",", "", ";", "", ".", "")
  for _, review := range reviews {
    // Copy over title
    curr_title := review.title
    // Format text
    curr_text := strings.ToLower(review.text)
    curr_text = replacer.Replace(curr_text)
    // Filter words out
    filterWords(&curr_text, filter_filename)

    formatted_text := strings.Fields(curr_text)
    // Loop through each word in text and input into index
    for i, word := range formatted_text {
      counted := false
      // Check to see if word is alredy in index
      _, in_index := index[word]
      _, in_counts := count_index[word]

      if !in_counts {
        count_index[word] = make(map[string] int)
      }
      _, title_in_counts := count_index[word][curr_title]

      if title_in_counts && !counted || !title_in_counts {
        count_index[word][curr_title] += 1
        counted = true
      }

      if !in_index {
        index[word] = make(map[string] []int)
      }
      index[word][curr_title] = append(index[word][curr_title], i)
    }

  }
  return index, count_index
}

func filterWords(text *string, filter_filename string) {
  file, err := os.Open(filter_filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  stop_words := readFile(file)

  // Loop through each of the stop words and filter
  for _, word := range stop_words {
    bound := "\b" // This little guy makes all the difference
    replacer := strings.NewReplacer(bound + word + bound, "")
    *text = replacer.Replace(*text)
  }
  // fmt.Println(*text)

}

func loadReviews(urls []string, num_urls int) []httpResponse {
  // the buffered channel that will block at num_urls
  buffered_chan := make(chan struct{}, num_urls)
  // the results channel that won't block
  results_chan := make(chan *httpResponse)

  defer func() {
    close(buffered_chan)
    close(results_chan)
  }()

  for _, url := range urls {
    go makeRequest(url, buffered_chan, results_chan)
  }

  var results []httpResponse

  for{
    result := <-results_chan
    results = append(results, *result)

    if len(results) == len(urls) {
  			break
  		}
  }
  return results
}


func readFile(file *os.File) []string {
  scanner := bufio.NewScanner(file)
  var tokens []string

  for scanner.Scan() {
    tokens = append(tokens, scanner.Text())
  }
  err := scanner.Err()
  if err != nil {
    log.Fatal(err)
  }
  return tokens
}
