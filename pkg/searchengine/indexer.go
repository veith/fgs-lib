package searchengine

import (
	"github.com/spf13/viper"
	"os"
)

// Adds an entry to the index or reindexing an existing entry
func Index(indexName string, id string, data interface{}) {
	cwd, _ := os.Getwd()
	p := cwd + viper.GetString("storage.basedir") + "/indexes/" + indexName

	index, _ := openIndex(p)
	index.Index(id, data)

}

// deletes an entry from the index
func DeleteFromIndex(indexName string, id string) {
	cwd, _ := os.Getwd()
	p := cwd + viper.GetString("storage.basedir") + "/indexes/" + indexName

	index, _ := openIndex(p)

	index.Delete(id)

}
