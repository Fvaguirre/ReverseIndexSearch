package main

import (
  "fmt"
  "bufio"
  "strings"
  "os"
  "./indexer"
)
var new_index indexer.Index

func init() {
  new_index = indexer.NewIndex("urls.txt", "stopWords.txt")
}

func main() {
  runSearchEngine()

}

func runSearchEngine() {
  // Loop forever asking for new search queries
  for {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("\nEnter 'q' to quit")
    fmt.Print("\nEnter a keyword to search: ")

    word, _ := reader.ReadString('\n')
    word = strings.TrimSpace(word)
    word = strings.ToLower(word)
    if word == "q" {
      break
    }
    if word == "" {
      fmt.Println("Error: Enter a keyword")
      continue
    }
    err := new_index.QueryWord(word)
    // If search query returns an error simply print it
    if err != nil {
      fmt.Println(err)
      break
    }
  }

}
