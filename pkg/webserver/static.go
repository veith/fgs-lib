package webserver

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
)

func Start(publicDir http.FileSystem) {

	http.Handle("/", http.StripPrefix("/", FileServerWith404(publicDir, fileSystem404)))

	target, _ := url.Parse(viper.GetString("server.transcoder.addr"))
	if target == nil {
		// no host given
		target, _ = url.Parse("http://localhost" + viper.GetString("server.transcoder.addr"))
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	http.HandleFunc("/api/", stripApiFromUrl(proxy))

	fmt.Printf("Webserver started on:" + viper.GetString("server.webserver.addr") + "\n")
	if err := http.ListenAndServe(viper.GetString("server.webserver.addr"), nil); err != nil {
		log.Fatal(err)
	}
}

func stripApiFromUrl(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// strip /api from url
		r.URL.Path = r.URL.Path[4:]
		r.Header.Set("api-base-url", "/api")
		p.ServeHTTP(w, r)
	}
}

// From https://gist.github.com/lummie/91cd1c18b2e32fa9f316862221a6fd5c
// FSHandler404 provides the function signature for passing to the FileServerWith404
type FSHandler404 = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

func fileSystem404(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool) {
	//if not found redirect to main index file for deeplinking...
	r.URL.Path = "/index.html"
	return true
}

func FileServerWith404(root http.FileSystem, handler404 FSHandler404) http.Handler {
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//make sure the url path starts with /
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		// attempt to open the file via the http.FileSystem
		f, err := root.Open(upath)
		if err != nil {
			if os.IsNotExist(err) {
				// call handler
				if handler404 != nil {
					doDefault := handler404(w, r)
					if !doDefault {
						return
					}
				}
			}
		}

		// close if successfully opened
		if err == nil {
			f.Close()
		}

		// default serve
		fs.ServeHTTP(w, r)
	})
}
