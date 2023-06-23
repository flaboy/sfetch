package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var pFlag = flag.Int("p", 6000, "Port listen")
var hostFlag = flag.String("h", "*", "host allow. [,] split")
var hostMap = map[string]bool{}

func main() {
	flag.Parse()

	log.Printf("listen port: %d\n", *pFlag)
	log.Printf("host allow: %s\n", *hostFlag)
	for _, host := range strings.Split(*hostFlag, ",") {
		hostMap[host] = true
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
	r2, err := httpClient.Do(req)
	if err != nil {
		rsp.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r2.Body.Close()
	for k, v := range r2.Header {
		rsp.Header()[k] = v
	}
	statusCode = r2.StatusCode
	rsp.WriteHeader(r2.StatusCode)
	io.Copy(rsp, r2.Body)
}
