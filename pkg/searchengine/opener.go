package searchengine

import (
	"github.com/blevesearch/bleve"
	"time"
)

var openIndexes = map[string]bleve.Index{}

func openIndex(path string) (bleve.Index, error) {
	if openIndexes[path] == nil {
		index, err := bleve.Open(path)
		if err != nil {
			// index does not exist
			createIndex(path)
			index, err = bleve.Open(path)
			if err != nil {
				return nil, err
			}
		}
		openIndexes[path] = index
		// close index after 60 minutes. Reopening costs us 20ms. So this will not hurt
		// so much and is easier then building a index session management
		go func() {
			time.Sleep(time.Minute * 60)
			index.Close()
			openIndexes[path] = nil
		}()
	}

	return openIndexes[path], nil
}

func createIndex(path string) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()

	index, err := bleve.New(path, mapping)
	defer index.Close()
	return index, err
}
