// Package 'indexer' provides the internals for building a map-based search
// index; It provides the needed file handling, http requests and map (index)
// creation
package indexer

import(
  "fmt"
  "strings"
  "errors"
)

// Index represents the word index; the index cleanses the text of punctuations
// {'.', ';', ',', '!'} and maps a given word to a college name (title) and then
// maps said name to a slice of ints representing where in the review body the
// word was found
type Index struct {
  index map[string]map[string][]int // Maps word -> title -> [] indexes
  filtered map[string]bool // Maps whether a word has been filtered
}

// Instantiates a new index object.
// Params: two string filenames containing the names of the list of urls and
// words to filter
// Returns: a new Index object populated with corresponding index and filter
// word map
func NewIndex(urls_file string, filter_file string) (Index) {
  var newIndex Index
  newIndex.index,newIndex.filtered = runIndexer(urls_file, filter_file)
  return newIndex
}

// Queries the index for a given word.
// Params: a string (word)
// Returns: an error if the word is not alphanumeric, or if more than one word,
// or if search word is a filtered word
func (i Index) QueryWord(word string) error {
  is_valid := isAlphaNumeric(word)
  if split := strings.Split(word, " "); len(split) > 1 {
    return errors.New("Error: Enter a single keyword")
  }
  if !is_valid  {
    return errors.New("Error: Enter a valid alphanumeric keyword")
  }
  if isFilteredWord(word, i.filtered) {
    return errors.New("Error: Enter a different keyword")
  }
  // Get our sorted int keys and our map
  sorted_keys, m := indexSearchWord(word, i.index)
  // Print out what we found!
  for _, key := range sorted_keys {
    for _, title := range m[key] {
      fmt.Printf("In %d reviews for %s\n", key, title)

    }
  }
  return nil
}
