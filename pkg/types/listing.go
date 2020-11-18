package types

import (
	"strconv"
	"strings"
)

// use this type for listings with options
type ListingOptions struct {
	//Partial representation (comma separated list of field names), ?fields=
	Fields string `json:"fields,omitempty"`
	//The response message will be filtered by the fields before being sent back to the client, filter=[[&#39;id&#39;,&#39;eq&#39;,&#39;1&#39;]]
	Filter string `json:"filter,omitempty"`
	//Specifies the result ordering for List requests. The default sorting order is ascending, ?order_by=foo desc,bar
	OrderBy string `json:"order_by,omitempty"`
	//Use this field to specify the maximum number of results to be returned by the server.
	//  The server may further constrain the maximum number of results returned in a single page.
	//  If the page_size is 0, the server will decide the number of results to be returned. page_size=15
	PageSize uint32 `json:"page_size,omitempty"`
	// The page to display
	Page uint32 `json:"page,omitempty"`
	// The cursor to display a page
	// use this if you have a lot of new data while paginating...
	PageCursor string `json:"page_cursor,omitempty"`
	//Query term to search
	Q string `json:"q,omitempty"`
	//allows the client to specify which view of the resource it wants to receive in the response. view=menu
	View string `json:"view,omitempty"`
}

// use this as response for lists with options
type ListingMetas struct {
	// the total count of records for the current request
	NumOfRecordsForRequest uint32
	// the page size which was used must not be the same like the szize requested
	UsedPageSize uint32 `json:"used_page_size,omitempty"`
}

func ConvertToURLQuery(l ListingOptions) string {
	//todo: implement ConvertToURLQuery
	var query []string
	if l.Fields != "" {
		query = append(query, "fields="+l.Fields)
	}

	if l.Filter != "" {
		query = append(query, "filter="+l.Filter)
	}

	if l.OrderBy != "" {
		query = append(query, "order_by="+l.OrderBy)
	}

	if l.PageSize != 0 {
		query = append(query, "page_size="+strconv.FormatUint(uint64(l.PageSize), 10))
	}
	if l.Page != 0 {
		query = append(query, "page="+strconv.FormatUint(uint64(l.Page), 10))
	}

	if l.Q != "" {
		query = append(query, "q="+l.Q)
	}

	if l.View != "" {
		query = append(query, "view="+l.View)
	}

	if len(query) > 0 {
		return "?" + strings.Join(query, "&")
	}
	return ""
}
