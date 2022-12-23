package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func MyNewSingleHostReverseProxy(target *url.URL, rewriteHost bool) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		if rewriteHost {
			req.Host = target.Host
		}
		req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string, rewriteHost bool) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return MyNewSingleHostReverseProxy(url, rewriteHost), nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

type ForwardCfg struct {
	Path        string
	TargetURL   string
	rewriteHost bool
}

func main() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})

	// simple forward
	cfgs := []ForwardCfg{
		{"/blog", "https://example.com", false},
		{"/", "https://tianpan.co", true},
	}
	for _, m := range cfgs {
		site, err := NewProxy(m.TargetURL, m.rewriteHost)
		if err != nil {
			panic(err)
		}
		rtr.PathPrefix(m.Path).HandlerFunc(ProxyRequestHandler(site))
	}

	http.Handle("/", rtr)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4321"
	}
	fmt.Printf("HTTP server listening on http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
