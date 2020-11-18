package pagination

import (
	"github.com/spf13/viper"
	"github.com/veith/fgs-lib/pkg/types"
)

var DefaultPageSize = uint32(23) //yes, as string

func GetListingOptions(generalRequest interface{}) types.ListingOptions {
	defaultOptioins := types.ListingOptions{
		Fields:     "",
		Filter:     "",
		OrderBy:    "",
		PageSize:   DefaultPageSize,
		Page:       0,
		PageCursor: "0",
		Q:          "",
		View:       "",
	}

	fld, hasFields := generalRequest.(Fields)
	if hasFields {
		defaultOptioins.Fields = fld.GetFields()
	}

	f, hasFilter := generalRequest.(Filter)
	if hasFilter {
		defaultOptioins.Filter = f.GetFilter()
	}

	ob, hasOrderBy := generalRequest.(OrderBy)
	if hasOrderBy {
		defaultOptioins.OrderBy = ob.GetOrderBy()
	}

	ps, hasPageSize := generalRequest.(PageSize)
	if hasPageSize {
		maxPageSize := viper.GetUint32("modules.all.max_items_per_page")

		defaultOptioins.PageSize = ps.GetPageSize()

		requestedPageSize := defaultOptioins.PageSize
		// use default on error
		if requestedPageSize == 0 {
			defaultOptioins.PageSize = DefaultPageSize
		}
		// set to max on exeeded values
		if requestedPageSize > maxPageSize {
			defaultOptioins.PageSize = maxPageSize
		}
	}

	p, hasPage := generalRequest.(Page)
	if hasPage {
		defaultOptioins.Page = p.GetPage()
	}

	pc, hasPageCursor := generalRequest.(PageCursor)
	if hasPageCursor {
		defaultOptioins.PageCursor = pc.GetPageCursor()
	}

	q, hasQ := generalRequest.(Q)
	if hasQ {
		defaultOptioins.Q = q.GetQ()
	}

	v, hasView := generalRequest.(View)
	if hasView {
		defaultOptioins.View = v.GetView()
	}

	return defaultOptioins
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
