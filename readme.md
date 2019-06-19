### Niche Review Search Engine
This program builds an index from a two given txt files of "\n" separated urls
and words to filter. Index represents the word index; the index cleanses the text
of {'.', ';', ',', '!'} and maps a given word to a college name (title) and then
maps said name to a slice of ints representing where in the review body the word
was found!

#### The indexer Package
The indexer package loads the data into memory and builds the given index! The
package comes with the exportable **Index** class and its two functions
**indexer.NewIndex() (Index)**, and **indexer.QueryWord(word, string) (error)**
which operates on an Index object!

### Running the Search Engine
* Run reviewSearchEngine.go with run reviewSearchEngine.go
* Enter a valid keyword when prompted
* Press q and enter to quit
