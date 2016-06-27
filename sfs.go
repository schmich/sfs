package main

import (
  "os"
  "fmt"
  "strconv"
  "net/http"
  "path/filepath"
  "github.com/jawher/mow.cli"
)

func main() {
  app := cli.App("sfs", "Static file server - https://github.com/schmich/sfs")
  app.Version("v version", "sfs " + Version)

  port := app.IntOpt("p port", 8080, "Listening port")
  iface := app.StringOpt("i iface interface", "127.0.0.1", "Listening interface")
  dir := app.StringOpt("d dir directory", ".", "Root directory to serve")

  var err error

  app.Action = func () {
    *dir, err = filepath.Abs(*dir)
    if err != nil {
      panic(err)
    }

    listen := *iface + ":" + strconv.Itoa(*port)

    fmt.Printf(">> Serving %s\n", *dir)
    fmt.Printf(">> Listening on %s\n", listen)
    fmt.Println(">> Ctrl+C to stop")

    server := http.FileServer(http.Dir(*dir))
    panic(http.ListenAndServe(listen, server))
  }

  app.Run(os.Args)
}
