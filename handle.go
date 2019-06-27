package main

import (
    "crypto/hmac"
    "crypto/sha1"
    "fmt"
    "github.com/gorilla/mux"
    "io"
    "log"
    "net/http"
    "net/url"
)

func digestHandler(w http.ResponseWriter, req *http.Request) {

    path := mux.Vars(req)
    query, err := url.ParseQuery(req.URL.RawQuery)

    if err != nil {
        log.Fatal(err.Error())
    }

    if path["digest"] != "" && query["url"][0] != "" {

        destUrl, err := url.QueryUnescape(query["url"][0])
        if err != nil {
            log.Fatal(err.Error())
        }

        hmc := hmac.New(sha1.New, []byte(sharedKey))
        io.WriteString(hmc, destUrl)
        sum := fmt.Sprintf("%x",hmc.Sum(nil))

        if path["digest"] == sum {
            destUrl, err := url.Parse(destUrl)

            if err != nil {
                log.Fatal(err.Error())
            }

            proxyImageRequest(&w, destUrl)
        } else {
            w.WriteHeader(404)
            io.WriteString(w, "not found")
        }

    } else {
        w.WriteHeader(404)
        io.WriteString(w, "not found")
    }
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "pong")
}

func statusHandler(w http.ResponseWriter, _ *http.Request) {
    io.WriteString(w, "status called")
}