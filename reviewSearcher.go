package main

import (
  "fmt"
  "./indexer"
)
var new_index indexer.Index

func init() {
  new_index = indexer.NewIndex("urls.txt", "stopWords.txt")
}

func main() {
  // Loop forever asking for new search queries
  for {
    err := new_index.QueryWords()
    // If search query returns an error simply print it
    if err != nil {
      fmt.Println(err)
    }
  }

}
