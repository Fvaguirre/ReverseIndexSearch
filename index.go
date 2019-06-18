package indexer

import (
  "./index_builder"
)


type index struct {
  index map[string]map[string][]int
  counts map[string]map[string]int
}

func NewIndex(urls_file string, stopWords_file string) (index) {
  var newIndex index
  newIndex.index, newIndex.counts = runIndexer("urls.txt", "stop_words.txt")

  return newIndex

}
