package searchengine

import (
	"github.com/blevesearch/bleve"
	"github.com/spf13/viper"
	"github.com/veith/fgs-lib/pkg/types"

	"os"
	"strings"
)

// list and filter collections,...
func List(indexName string, options types.ListingOptions, indexes map[string]string) (*bleve.SearchResult, error) {
	cwd, _ := os.Getwd()
	p := cwd + viper.GetString("storage.basedir") + "/indexes/" + indexName

	index, err := openIndex(p)

	items_by_page := options.PageSize
	if items_by_page == 0 {
		items_by_page = viper.GetUint32("modules.all.items_per_page")
	}

	// bleve start with page 1
	page := options.Page
	if page <= 0 {
		page = 1
	} else {
		page++
	}

	// field scoping
	fields := []string{}
	for k, v := range indexes {
		fields = append(fields, k+":"+v)
	}
	fieldscope := strings.Join(fields, " ")

	var searchRequest *bleve.SearchRequest

	if options.Q == "" {

		bq := bleve.NewBooleanQuery()
		bq.Must = bleve.NewMatchQuery(fieldscope)
		bq.Should = bleve.NewMatchAllQuery()
		query := bq
		searchRequest = bleve.NewSearchRequestOptions(query, int(items_by_page), int((page-1)*items_by_page), false)
	} else {
		query := bleve.NewFuzzyQuery(fieldscope + options.Q)
		searchRequest = bleve.NewSearchRequestOptions(query, int(items_by_page), int((page-1)*items_by_page), false)
	}
	// default sort order is id desc
	sortOrder := []string{"-_id"}

	if options.OrderBy != "" {
		sortOrder = strings.Split(strings.ReplaceAll(options.OrderBy, " ", ""), ",")
	}

	//searchRequest.Fields = strings.Split(strings.ReplaceAll(options.Fields, " ", ""), ",")
	// todo implement options.Filter

	searchRequest.SortBy(sortOrder)

	res, err := index.Search(searchRequest)
	return res, err
}
