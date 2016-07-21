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

func main() {
  app := cli.App("sfs", "Static file server - https://github.com/schmich/sfs")
  app.Version("v version", "sfs " + Version)

  port := app.IntOpt("p port", 8080, "Listening port")
  iface := app.StringOpt("i iface interface", "127.0.0.1", "Listening interface")
  dir := app.StringOpt("d dir directory", ".", "Directory to serve")
  allIface := app.BoolOpt("g global", false, "Listen on all interfaces (overrides -i)")
  noBrowser := app.BoolOpt("B no-browser", false, "Do not launch browser")

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

    fmt.Printf(">> Serving %s\n", *dir)
    fmt.Printf(">> Listening on %s\n", listen)
    fmt.Println(">> Ctrl+C to stop")

    if !*noBrowser {
      url := "http://127.0.0.1" + portPart
      open.Start(url)
    }

    server := http.FileServer(http.Dir(*dir))
    panic(http.ListenAndServe(listen, server))
  }

  app.Run(os.Args)
}
