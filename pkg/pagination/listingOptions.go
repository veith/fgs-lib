package pagination

import (
	"github.com/spf13/viper"
	"github.com/veith/fgs-lib/pkg/types"
)

var DefaultPageSize = uint32(23) //yes, as string

func GetListingOptions(generalRequest interface{}) types.ListingOptions {
	defaultOptions := types.ListingOptions{
		Fields:     "",
		Filter:     "",
		OrderBy:    "",
		PageSize:   DefaultPageSize,
		Page:       0,
		PageCursor: "",
		Q:          "",
		View:       "",
	}

	fld, hasFields := generalRequest.(Fields)
	if hasFields {
		defaultOptions.Fields = fld.GetFields()
	}

	f, hasFilter := generalRequest.(Filter)
	if hasFilter {
		defaultOptions.Filter = f.GetFilter()
	}

	ob, hasOrderBy := generalRequest.(OrderBy)
	if hasOrderBy {
		defaultOptions.OrderBy = ob.GetOrderBy()
	}

	ps, hasPageSize := generalRequest.(PageSize)
	if hasPageSize {
		maxPageSize := viper.GetUint32("modules.all.max_items_per_page")

		defaultOptions.PageSize = ps.GetPageSize()

		requestedPageSize := defaultOptions.PageSize
		// use default on error
		if requestedPageSize == 0 {
			defaultOptions.PageSize = DefaultPageSize
		}
		// set to max on exeeded values
		if requestedPageSize > maxPageSize {
			defaultOptions.PageSize = maxPageSize
		}
	}

	p, hasPage := generalRequest.(Page)
	if hasPage {
		defaultOptions.Page = p.GetPage()
	}

	pc, hasPageCursor := generalRequest.(PageCursor)
	if hasPageCursor {
		defaultOptions.PageCursor = pc.GetPageCursor()
	}

	q, hasQ := generalRequest.(Q)
	if hasQ {
		defaultOptions.Q = q.GetQ()
	}

	v, hasView := generalRequest.(View)
	if hasView {
		defaultOptions.View = v.GetView()
	}

	return defaultOptions
}

type Fields interface {
	GetFields() string
}

type Filter interface {
	GetFilter() string
}

type OrderBy interface {
	GetOrderBy() string
}
type PageSize interface {
	GetPageSize() uint32
}
type PageCursor interface {
	GetPageCursor() string
}
type View interface {
	GetView() string
}

type Page interface {
	GetPage() uint32
}
type Q interface {
	GetQ() string
}
