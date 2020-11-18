package hateoas

import (
	"context"
	"github.com/spf13/viper"
	furopb "github.com/theNorstroem/FuroBaseSpecs/dist/pb/furo"
	"github.com/veith/fgs-lib/pkg/types"
	"google.golang.org/grpc/metadata"
	"strings"
)

type Builder struct {
	ServiceName string
}

type Idmap map[string]string

func NewHTSBuilder(ServiceName string) *Builder {
	return &Builder{
		ServiceName: ServiceName,
	}
}

/**
 * Reads the first api-base-url from context
 * The first is the one which is outermost.
 */
func BaseUrlFromContext(ctx context.Context) string {
	// Context api-base-url is provided from grpc-server which get it from grpc-gateway
	// we use it to have absolute hrefs for the HATEOAS
	// if api-base-url header was not set, we provide relative paths to the client
	md, _ := metadata.FromIncomingContext(ctx)
	baseURL := "/"
	if len(md["api-base-url"]) > 0 {
		baseURL = md["api-base-url"][0]
	}
	return baseURL
}

func (b *Builder) EntityHTS(ctx context.Context, entitytype string, ids Idmap, urlPattern string, htsMethods ...string) []*furopb.Link {
	links := []*furopb.Link{}

	baseURL := BaseUrlFromContext(ctx)

	links = append(links, &furopb.Link{
		Href:    baseURL + ReplacePaternWithIdmap(urlPattern, ids),
		Method:  "GET",
		Rel:     "self",
		Service: b.ServiceName,
		Type:    entitytype,
	})

	// Append HTS for other then self
	for _, method := range htsMethods {
		link := &furopb.Link{
			Href:    baseURL + ReplacePaternWithIdmap(urlPattern, ids),
			Method:  strings.ToUpper(method),
			Service: b.ServiceName,
			Type:    entitytype,
		}
		switch method {
		case "PATCH":
			link.Rel = "update"
			break
		case "PUT":
			link.Rel = "update"
			break
		case "DELETE":
			link.Rel = "delete"
			break
		default:
			link.Rel = strings.ToLower(method)
		}

		links = append(links, link)
	}

	return links
}

func ReplacePaternWithIdmap(pattern string, idmap Idmap) string {
	for k, v := range idmap {
		pattern = strings.Replace(pattern, "{"+k+"}", v, -1)
	}
	return pattern
}

func (b *Builder) CollectionHTS(ctx context.Context, listingOptions types.ListingOptions, listingMetas types.ListingMetas, idmap Idmap, urlPattern string, s ...string) []*furopb.Link {
	links := []*furopb.Link{}
	// Context api-base-url is provided from grpc-server which get it from grpc-gateway
	// we use it to have absolute hrefs for the HATEOAS
	// if api-base-url header was not set, we provide relative paths to the client
	md, _ := metadata.FromIncomingContext(ctx)
	// fallback config
	baseURL := "/"
	if len(md["api-base-url"]) > 0 {
		baseURL = md["api-base-url"][0]
	}

	listingOptions.PageSize = listingMetas.UsedPageSize

	// self
	AddLinkToLinksArray(&links, &furopb.Link{
		Href:    baseURL + ReplacePaternWithIdmap(urlPattern, idmap) + types.ConvertToURLQuery(listingOptions),
		Rel:     "self",
		Service: b.ServiceName,
	}, "")
	// self
	AddLinkToLinksArray(&links, &furopb.Link{
		Href:    baseURL + ReplacePaternWithIdmap(urlPattern, idmap) + types.ConvertToURLQuery(listingOptions),
		Rel:     "list",
		Service: b.ServiceName,
	}, "")

	// create
	AddLinkToLinksArray(&links, &furopb.Link{
		Href:    baseURL + ReplacePaternWithIdmap(urlPattern, idmap) + types.ConvertToURLQuery(listingOptions),
		Method:  "POST",
		Rel:     "create",
		Service: b.ServiceName,
	}, "")

	// if we are on page > 1 we add a prev
	currentPage := listingOptions.Page
	lo := listingOptions
	pageSize := listingOptions.PageSize

	if listingMetas.UsedPageSize > 0 {
		pageSize = listingMetas.UsedPageSize
	}
	if pageSize == 0 {
		pageSize = viper.GetUint32("modules.all.items_per_page")
	}
	lastPage := listingMetas.NumOfRecordsForRequest / pageSize

	lo.Page = lastPage
	AddLinkToLinksArray(&links, &furopb.Link{
		Href:    baseURL + ReplacePaternWithIdmap(urlPattern, idmap) + types.ConvertToURLQuery(lo),
		Rel:     "last",
		Service: b.ServiceName,
	}, "")

	// first page is 1
	lo.Page = 0
	AddLinkToLinksArray(&links, &furopb.Link{
		Href:    baseURL + ReplacePaternWithIdmap(urlPattern, idmap) + types.ConvertToURLQuery(lo),
		Rel:     "first",
		Service: b.ServiceName,
	}, "")

	// only give a prev if we are on the first site
	if currentPage > 0 {
		// prev
		lo.Page = currentPage - 1
		AddLinkToLinksArray(&links, &furopb.Link{
			Href:    baseURL + ReplacePaternWithIdmap(urlPattern, idmap) + types.ConvertToURLQuery(lo),
			Rel:     "prev",
			Service: b.ServiceName,
		}, "")
	}

	// next
	if currentPage < lastPage {
		lo.Page = currentPage + 1
		AddLinkToLinksArray(&links, &furopb.Link{
			Href:    baseURL + ReplacePaternWithIdmap(urlPattern, idmap) + types.ConvertToURLQuery(lo),
			Rel:     "next",
			Service: b.ServiceName,
		}, "")
	}

	return links
}

// Add a HTS link. If Method is not given, GET is used
func AddLinkToLinksArray(target *[]*furopb.Link, link *furopb.Link, linkType string) {
	link.Type = linkType
	// make method GET the default
	if link.Method == "" {
		link.Method = "GET"
	}
	// make self the default
	if link.Rel == "" {
		link.Rel = "self"
	}
	*target = append(*target, link)
}
