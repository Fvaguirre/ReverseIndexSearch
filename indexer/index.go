package indexer

import(
  "fmt"
  "bufio"
  "os"
  "strings"
  "errors"
)

type Index struct {
  index map[string]map[string][]int
  counts map[string]map[string]int
  filtered map[string]bool
}

func NewIndex(urls_file string, filter_file string) (Index) {
  var newIndex Index
  newIndex.index, newIndex.counts, newIndex.filtered = runIndexer(urls_file, filter_file)
  return newIndex
}



func (i Index) QueryWords() error {
  reader := bufio.NewReader(os.Stdin)

  fmt.Print("\nEnter a keyword to search: ")

  words, _ := reader.ReadString('\n')
  words = strings.TrimSpace(words)
  words = strings.ToLower(words)
  split_words := strings.Split(words, " ")
  is_valid := isValidQuery(split_words, i.filtered)
  if split_words[0] == "" {
    return errors.New("Error: Enter a keyword")
  }
  if !is_valid  {
    return errors.New("Error: Enter a valid keyword or keyword phrase")
  }
  if len(split_words) < 2 && isFilteredWord(split_words[0], i.filtered) {
    return errors.New("Error: Enter a different keyword")
  }
  if len(split_words) > 1 {
    return errors.New("Error: Enter a single keyword")
  }

  sorted_keys, m := indexSearchWord(words, i.counts)
  fmt.Println(sorted_keys)
  if len(sorted_keys) == 0 {
    return errors.New("No search results")
  }
  for _, key := range sorted_keys {
    for _, title := range m[key] {
      fmt.Printf("In %d reviews for %s\n", key, title)

    }
  }

  return nil
}
