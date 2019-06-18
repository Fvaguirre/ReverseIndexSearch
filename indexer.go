package main

import(
  "bufio"
  "io/ioutil"
  "fmt"
  "os"
  "net/http"
  "time"
  "log"
  "strings"
)

type httpResponse struct {
  title string
  text string
  err error
}

func main(){
  // Try to open text file with urls
  file, err := os.Open("urls.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  // urls := readUrls(file)
  runIndexer("urls.txt", "stopWords.txt")
  // fmt.Println(urls[0])
  // start := time.Now()
  // ch := make(chan string)
  // Loop through input file

}
func runIndexer(filename string, filename2 string) {
  // Try to open text file with urls
  file, err := os.Open(filename)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  urls := readFile(file)
  fmt.Println(len(urls))
  // fmt.Println(urls[0])
  // num_urls := len(urls)
  start := time.Now()
  res := loadReviews(urls, 100)
  end := time.Since(start)
  fmt.Printf("Time to loadReviews into memory: %s", end)
  fmt.Println(len(res))

  start = time.Now()
  _, counts_index := indexReviews(res, filename2)
  end = time.Since(start)
  fmt.Printf("Time to build index: %s",end )
  fmt.Println(len(counts_index))
  // for _, val:= range counts_index["study"] {
  //   fmt.Printf(" val: %d\n", val)
  // }
  // fmt.Println(counts_index["study"])
  // for key, _ := range index {
  //   fmt.Printf("Key:%s has value: ", key)
    // for key2, val2 := range val {
    //   fmt.Printf("[%s] = [", key2)
    //
    //   for _, val3 := range val2 {
    //     fmt.Printf("%d, ", val3)
    //   }
    //   fmt.Printf("]\n")
    // }
  // }
  res = nil

  // fmt.Println(index)
  // fmt.Println(index["University of Phoenix - Dallas"])
  // fmt.Println(res)
  // size := len(res)
  // fmt.Println(size)
  //
  // fmt.Println(res[0].title)
}
func indexReviews(reviews []httpResponse, filter_words string) (map[string] map[string] []int, map[string] map[string] int) {
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
    filterWords(&curr_text, filter_words)

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

func filterWords(text *string, file_name string) {
  file, err := os.Open(file_name)
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
  // count := 0
  for scanner.Scan() {
    // count += 1
    tokens = append(tokens, scanner.Text())
  }
  err := scanner.Err()
  if err != nil {
    log.Fatal(err)
  }
  // fmt.Println(count)
  return tokens
}
// func makeRequest(url string, url_map sync.Map) {
//   // fmt.Println("AHHAA")
//   // start := time.Now()
//   resp, _  := http.Get(url)
//   // secs := time.Since(start).Seconds()
//   body, _ := ioutil.ReadAll(resp.Body)
//   // fmt.Println(body)
//   url_map.Store(url, body)
//   // str := fmt.Sprintf("%.2f elapsed with response length %d %s", secs, len(body), url)
//   // fmt.Println(str)
// }
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
