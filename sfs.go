package main

import (
  "os"
  "fmt"
  "strconv"
  "net/http"
  "path/filepath"
  "github.com/jawher/mow.cli"
  "github.com/skratchdot/open-golang/open"
)

func NoCache(h http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Cache-Control", "private, max-age=0, no-cache")
    r.Header.Set("Cache-Control", "no-cache")
    r.Header.Del("If-Modified-Since")
    r.Header.Del("If-None-Match")
    h.ServeHTTP(w, r)
  })
}

func main() {
  app := cli.App("sfs", "Static file server - https://github.com/schmich/sfs")
  app.Version("v version", "sfs " + Version)

  port := app.IntOpt("p port", 8080, "Listening port")
  iface := app.StringOpt("i iface interface", "127.0.0.1", "Listening interface")
  allIface := app.BoolOpt("g global", false, "Listen on all interfaces (overrides -i)")
  dir := app.StringOpt("d dir directory", ".", "Directory to serve")
  noBrowser := app.BoolOpt("B no-browser", false, "Do not launch browser")
  cache := app.BoolOpt("c cache", false, "Allow cached responses")

  app.Action = func () {
    var err error

    *dir, err = filepath.Abs(*dir)
    if err != nil {
      panic(err)
    }

    if *allIface {
      *iface = "0.0.0.0"
    }

    portPart := ":" + strconv.Itoa(*port)
    listen := *iface + portPart

    server := http.FileServer(http.Dir(*dir))
    if !*cache {
      server = NoCache(server)
    }

    fmt.Printf(">> Serving %s\n", *dir)
    fmt.Printf(">> Listening on %s\n", listen)
    fmt.Println(">> Ctrl+C to stop")

    if !*noBrowser {
      url := "http://127.0.0.1" + portPart
      open.Start(url)
    }

    panic(http.ListenAndServe(listen, server))
  }

  app.Run(os.Args)
}
