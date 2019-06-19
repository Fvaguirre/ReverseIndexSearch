// Package 'indexer' provides the internals for building a map-based search
// index; It provides the needed file handling, http requests and map (index)
// creation
package indexer

import (
  "sort"
)

// Checks whether given word is stored within given map.
// Params: word string, filtered map[string]bool
// Returns: bool; true if contained in map, false otherwise
func isFilteredWord(word string, filtered map[string]bool) bool{
  if _, in_map := filtered[word]; in_map {
    return true
  }
  return false
}

// Checks whether a given string is alphanumeric.
// Params: a string
// Returns: bool true if string is is alphanumeric; false otherwise
func isAlphaNumeric(s string) bool {
  for _, r := range s{
    if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
      continue
    }
    return false
  }
  return true
}

// Checks whether a given integer is present in an unsorted slice of ints.
// Params: an int value, and the slice of ints
// Returns true if int is already within slice, false otherwise
func isIn(val int, vals []int) bool{
  for _, v := range vals {
    if v == val {
      return true;
    }
  }
  return false
}

// Searches the map (index) for given word, then creates the inverted map and
// sorted key slice.
// Params: word (string) we are searching for and a map[string]map[string][]int
// to search in
// Returns: a sorted integer slice corresponding to keys in the returned
// map[int][]string which maps int counts to slices of strings containing titles 
func indexSearchWord(word string, index map[string]map[string][]int) ([]int, map[int][]string){
  // Build new map with int keys (counts) and colleges as vals
  inv_map := make(map[int][]string)
  // Keep int keys in separate arrays to sort
  var keys []int
  // Check if word is in index
  if title_indexes, ok := index[word]; ok {
    // Populate both inv_map and keys
    for key, val := range title_indexes {
      // Add title to int count key
      inv_map[len(val)] = append(inv_map[len(val)], key)
      if in := isIn(len(val), keys); !in {
        keys = append(keys, len(val))
      }
    }
    // Sort keys
    sort.Sort(sort.Reverse(sort.IntSlice(keys)))
    // Go through and sort the titles in order
    for _, val := range inv_map {
      sort.Strings(val)
    }
  }
  return keys, inv_map
}
