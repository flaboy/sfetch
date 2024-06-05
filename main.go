package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var pFlag = flag.Int("p", 6000, "Port listen")
var hostFlag = flag.String("h", "*", "host allow. [,] split")
var hostMap = map[string]bool{}
var rspHeaderMap = map[string]bool{}

func main() {
	flag.Parse()

	log.Printf("listen port: %d\n", *pFlag)
	log.Printf("host allow: %s\n", *hostFlag)
	for _, host := range strings.Split(*hostFlag, ",") {
		hostMap[host] = true
	}

	rspHeaderMap = map[string]bool{
		"Content-Type":        true,
		"Content-Length":      true,
		"Content-Encoding":    true,
		"Transfer-Encoding":   true,
		"Content-Disposition": true,
		"Date":                true,
		"Expires":             true,
		"Server":              true,
		"Vary":                true,
	}

	s := &http.Server{
		Addr:    ":" + strconv.Itoa(*pFlag),
		Handler: Handler{},
	}
	log.Fatal(s.ListenAndServe())
}

func checkHost(host string) bool {
	if len(hostMap) == 0 {
		return true
	}
	if _, ok := hostMap["*"]; ok {
		return true
	}
	_, ok := hostMap[host]
	return ok
}

type Handler struct {
}

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}
}

// make proxy /a.com/xxx to https://a.com/xxx
func (h Handler) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	statusCode := 0
	defer func() {
		log.Printf("%s %s %s %s %d\n", req.RemoteAddr, req.Method, req.URL.Host, req.URL.Path, statusCode)
	}()

	paths := strings.Split(req.URL.Path, "/")
	if len(paths) < 2 {
		rsp.WriteHeader(http.StatusBadRequest)
		return
	}

	if paths[1] == "g" {
		//解析querystring, 从querystring中获取url
		urlBase64 := req.URL.Query().Get("u")
		if urlBase64 == "" {
			fmt.Println("urlBase64 is empty")
			rsp.WriteHeader(http.StatusBadRequest)
			return
		}

		theUrl, err := base64.StdEncoding.DecodeString(urlBase64)
		if err != nil {
			fmt.Println("base64 decode error:", err, urlBase64)
			rsp.WriteHeader(http.StatusBadRequest)
			return
		}

		u2, err := url.Parse(string(theUrl))
		if err != nil {
			fmt.Println("url parse error:", err, string(theUrl))
			rsp.WriteHeader(http.StatusBadRequest)
			return
		}

		req.URL = u2
		req.Host = u2.Host
		req.RequestURI = ""
	} else {
		host := paths[1]
		if host == "" {
			rsp.WriteHeader(http.StatusBadRequest)
			return
		}
		if !checkHost(host) {
			rsp.WriteHeader(http.StatusForbidden)
			return
		}
		req.URL.Scheme = "https"
		req.URL.Host = host
		req.URL.Path = "/" + strings.Join(paths[2:], "/")
		req.Host = host
		req.RequestURI = ""
	}

	// proxyURL, err := url.Parse("http://127.0.0.1:49193")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// transport := &http.Transport{
	// 	Proxy: http.ProxyURL(proxyURL),
	// }

	// httpClient = &http.Client{
	// 	Transport: transport,
	// }

	r2, err := httpClient.Do(req)
	if err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r2.Body.Close()
	for k, v := range r2.Header {
		if _, ok := rspHeaderMap[k]; ok {
			rsp.Header()[k] = v
		}
	}

	statusCode = r2.StatusCode
	rsp.WriteHeader(r2.StatusCode)
	io.Copy(rsp, r2.Body)
}
