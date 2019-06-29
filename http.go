package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func proxyImageRequest(writer *http.ResponseWriter, destUrl *url.URL, remainRedirect int) {

	if destUrl.Scheme == "http" || destUrl.Scheme == "https" {

		newRequest, err := http.NewRequest("GET", destUrl.String(), nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		addTransferredHeaders(newRequest)

		timeout := time.Duration(10 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}

		response, err := client.Do(newRequest)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer response.Body.Close()

		if contentLengthLimit <= response.ContentLength {
			(*writer).WriteHeader(404)
			io.WriteString(*writer, "Content Length Exceeded")
			return
		}

		switch response.StatusCode {
		case 301, 302, 303, 307:

			if remainRedirect <= 0 {

				(*writer).WriteHeader(404)
				io.WriteString(*writer, "Exceeded max depth")

			} else if location := response.Header.Get("Location"); location != "" {

				(*writer).WriteHeader(404)
				io.WriteString(*writer, "Redirect with no location")

			} else {

				newUrl, err := url.Parse(response.Header.Get("Location"))
				if err != nil {
					log.Fatal(err.Error())
				}

				if newUrl.Scheme == "" {
					newUrl.Scheme = destUrl.Scheme
				}

				if newUrl.Host == "" {
					newUrl.Host = destUrl.Host
				}

				proxyImageRequest(writer, newUrl, remainRedirect-1)

			}
		case 304:

			(*writer).WriteHeader(304)
			io.WriteString(*writer, "Not modified")

		default:

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Println(err.Error())
			}
			setNewHeader(writer, response)
			(*writer).Write(body)

		}

	} else {

		(*writer).WriteHeader(400)
		io.WriteString(*writer, "Unknown protocol")

	}
}

func addTransferredHeaders(request *http.Request) {
	request.Header.Set("Via", "Gamo Asset Proxy")
	request.Header.Set("User-Agent", "Gamo Asset Proxy")
	request.Header.Set("Accept", "image/*")
	request.Header.Set("Accept-Encoding", "")
	request.Header.Set("X-Frame-Options", "deny")
	request.Header.Set("X-XSS-Protection", "1; mode=block")
	request.Header.Set("X-Content-Type-Options", "nosniff")
	request.Header.Set("Content-Security-Policy", "default-src 'none'; img-src data:; style-src 'unsafe-inline'")
}

func setNewHeader(writer *http.ResponseWriter, response *http.Response) {

	headers := []string{"Etag", "Expires", "Last-Modified", "Transfer-Encoding", "Content-Encoding", "Content-Length", "Content-Type"}
	for _, v := range headers {
		if headerValue := response.Header.Get(v); headerValue != "" {
			(*writer).Header().Set(v, headerValue)
		}
	}

	if cacheControl := response.Header.Get("Cache-Control"); cacheControl != "" {
		(*writer).Header().Set("Cache-Control", cacheControl)
	} else {
		(*writer).Header().Set("Cache-Control", "public, max-age=31536000")
	}

	(*writer).Header().Set("X-Frame-Options", "deny")
	(*writer).Header().Set("X-XSS-Protection", "1; mode=block")
	(*writer).Header().Set("X-Content-Type-Options", "nosniff")
	(*writer).Header().Set("Content-Security-Policy", "default-src 'none'; img-src data:; style-src 'unsafe-inline'")
	(*writer).Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

}
