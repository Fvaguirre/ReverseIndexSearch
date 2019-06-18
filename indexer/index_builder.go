// Package 'indexer' provides the internals for building a map-based search
// index; It provides the needed file handling, http requests and map (index)
// creation
package indexer

import (
  "bufio"
  "io/ioutil"
  "fmt"
  "os"
  "net/http"
  "log"
  "strings"
)

// Struct httpResponse models a review return from a get request to a given url.
type httpResponse struct {
  title string // The name of the university
  text string // The body of the review
  url string // The url of the given review
  err error // Any error returned from the get request

}

// Creates the index in the form of three maps.
// Params: two string corresponding to filename containing review urls to index
// and a filename containing words to filter for index creation
// Returns: the inverted index map[string]map[string][]int, counts index map[string]
// map[string]int and the filtered words map map[string]bool
func runIndexer(filename string, filter_filename string) (map[string] map[string] []int, map[string] map[string] int, map[string]bool) {
  // Try to open text file with urls
  file, err := os.Open("./data/" + filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  // Get the urls in a slice
  urls := readFile(file)
  // Make the necessary get requests for all the urls, limit 200 at a time
  res := loadReviews(urls, 200)
  // Get our built indexes, and filtered words as a slice
  index , counts_index, filtered_words := indexReviews(res, "./data/" + filter_filename)

  // Convert filtered words to slice for quicker access
  filtered_words_map := make(map[string]bool)
  for _, fw := range filtered_words {
    filtered_words_map[fw] = true
  }
  return index, counts_index, filtered_words_map
}

// Reads a text file line by line and returns a string slice containing contents.
// Params: a pointer to the text file
// Returns: a string slice containg contents of file
func readFile(file *os.File) []string {
  scanner := bufio.NewScanner(file)
  var tokens []string
  // Add each line to end of slice
  for scanner.Scan() {
    tokens = append(tokens, scanner.Text())
  }
  err := scanner.Err()
  if err != nil {
    log.Fatal(err)
  }
  return tokens
}

// Makes http get requests from a list of urls and loads them into memory.
// Params: a string slice of urls and an integer limit for buffering
// Returns: a slice of httpResponse objects populated with the data returned from
// the get requests on the given urls
func loadReviews(urls []string, limit int) []httpResponse {
  // the buffered channel that will block at limit
  buffered_chan := make(chan struct{}, limit)
  // the results channel that won't block
  results_chan := make(chan *httpResponse)

  // Let's not forget to close our channels
  defer func() {
    close(buffered_chan)
    close(results_chan)
  }()

  // Concurrently make the http get requests
  for _, url := range urls {
    go makeRequest(url, buffered_chan, results_chan)
  }

  // Our return slice
  var results []httpResponse
  // Populate the return slice
  for{
    result := <-results_chan
    results = append(results, *result)
    if len(results) == len(urls) {
  			break
  		}
    }
  return results
}

// Makes a get request for a target url and uses a buffered channel and results
// channel to do so concurrently
// Params: string url, chan struct{} to buffer the get requests, and a chan
// *httpResponse to hold the resulting responses
// Modifies: buffered chan, results_chan
func makeRequest(url string, buffered_chan chan struct{}, results_chan chan <- *httpResponse) {
  // Placeholder for buffered chan
  buffered_chan <- struct{}{}
  // Get the review
  resp, err := http.Get(url)

  // load review body
  body, _ := ioutil.ReadAll(resp.Body)

  // Always make sure to close your response!
  defer resp.Body.Close()

  // Load response into httpResponse object
  result := &httpResponse{"", string(body), url, err}
  // Split the text at first line break to get title
  splitted := strings.SplitN(result.text, "\n", 2)
  result.title = splitted[0]

  // Move result into results chan
  results_chan <- result

  // Remove one from buffered_chan
  <- buffered_chan
}


func indexReviews(reviews []httpResponse, filter_filename string) (map[string] map[string] []int, map[string] map[string] int, []string) {
  index := make(map[string] (map[string] []int))
  count_index := make(map[string] (map[string] int))
  replacer := strings.NewReplacer(",", "", ";", "", ".", "", "!", "")
  filtered_words := getFilteredWords(filter_filename)
  for _, review := range reviews {
    fmt.Println("indexing")
    fmt.Println(review.url)
    // Copy over title
    curr_title := review.title
    // Format text
    curr_text := strings.ToLower(review.text)
    curr_text = replacer.Replace(curr_text)
    // Filter words out
    filterWords(&curr_text, filtered_words)
    // Format resulting text into slice
    formatted_text := strings.Fields(curr_text)
    // Loop through each word in text and input into index
    for i, word := range formatted_text {
      counted := false
      // Check to see if word is alredy in index
      _, in_index := index[word]
      _, in_counts := count_index[word]
      // If word not in the index yet then set its value to an empty map
      if !in_counts {
        count_index[word] = make(map[string] int)
      }
      // Check to see if the current review title has been inputted into the
      // index at for the current word
      _, title_in_counts := count_index[word][curr_title]

      // If the title has been inputted, but has yet to be counted or
      if title_in_counts && !counted || !title_in_counts {
        count_index[word][curr_title] += 1
        counted = true
      }

      if !in_index {
        index[word] = make(map[string] []int)
      }
      index[word][curr_title] = append(index[word][curr_title], i)
    }
    fmt.Println("Finished.")
  }

  return index, count_index, filtered_words
}

func getFilteredWords(filter_filename string) []string {
  file, err := os.Open(filter_filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  stop_words := readFile(file)
  return stop_words
}

func filterWords(text *string, filtered_words []string) {
  // Loop through each of the stop words and filter
  for _, word := range filtered_words {
    bound := "\b" // This little guy makes all the difference
    replacer := strings.NewReplacer(bound + word + bound, "")
    *text = replacer.Replace(*text)
  }
}
