package indexer

import (
  "sort"
  // "fmt"
)

func isFilteredWord(word string, filtered map[string]bool) bool{
  if _, in_map := filtered[word]; in_map {
    return true
  }
  return false
}
func isValidQuery(words []string, filtered map[string]bool) bool{
  for _, word := range words {
    if !isAlphaNumeric(word) {
      return false
    }
  }
  return true
}
func isAlphaNumeric(s string) bool {
  for _, r := range s{
    if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
      continue
    }
    return false
  }
  return true
}

func binSearch(val int , vals []int , l int, r int) int{
  if r >= l {
    mid := l + (r-l)/2
    if vals[mid] == val {
      return mid
    }

    if vals[mid] > val {
      return binSearch(val, vals, l, mid - 1)
    }
    return binSearch(val, vals, mid + 1, r)
  }
  return -1
}
func indexSearchWord(word string, counts map[string]map[string]int) ([]int, map[int][]string) {
  if frequencies, ok := counts[word]; ok {

    // Build new map with int counts as key and colleges as values
    inv_map := make(map[int][]string)
    // Keep int keys in separate array to sort
    var keys []int
    for key, val := range frequencies {
      inv_map[val] = append(inv_map[val], key)
      // If no keys or current key hasnt been put in yet
      in := binSearch(val, keys, 0, len(keys)-1)
      if len(keys) == 0 || in == -1 {
        keys = append(keys, val)
      }
    }

    // Sort keys
    sort.Sort(sort.Reverse(sort.IntSlice(keys)))
    // Sort internal inv_map arrays
    for _, val := range inv_map {
      sort.Strings(val)
    }
    return keys, inv_map

  }
  return nil, nil

}
