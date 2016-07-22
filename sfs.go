package main

import (
  "os"
  "fmt"
  "time"
  "strconv"
  "strings"
  "regexp"
  "net/http"
  "path/filepath"
  "github.com/jawher/mow.cli"
  "github.com/skratchdot/open-golang/open"
)

type TraceResponseWriter struct {
  impl http.ResponseWriter
  bytesWritten int
  statusCode int
}

func NewTraceResponseWriter(writer http.ResponseWriter) *TraceResponseWriter {
  return &TraceResponseWriter{
    impl: writer,
    bytesWritten: 0,
    statusCode: 200,
  }
}

func (writer *TraceResponseWriter) Header() http.Header {
  return writer.impl.Header()
}

func (writer *TraceResponseWriter) Write(bytes []byte) (int, error) {
  writer.bytesWritten += len(bytes)
  return writer.impl.Write(bytes)
}

func (writer *TraceResponseWriter) WriteHeader(statusCode int) {
  writer.statusCode = statusCode
  writer.impl.WriteHeader(statusCode)
}

func formatSize(bytes int) string {
  if bytes < 1000 {
    return strconv.Itoa(bytes)
  } else if bytes < 1000000 {
    return fmt.Sprintf("%.2fK", float32(bytes) / 1000)
  } else {
    return fmt.Sprintf("%.2fM", float32(bytes) / 1000000)
  }
}

func TraceServer(h http.Handler, log string) http.Handler {
  formatter, _ := regexp.Compile("%.")

  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    traceWriter := NewTraceResponseWriter(w)
    h.ServeHTTP(traceWriter, r)

    ip := r.RemoteAddr
    if i := strings.Index(ip, ":"); i >= 0 {
      ip = ip[:i]
    }

    line := string(formatter.ReplaceAllFunc([]byte(log), func(match []byte) []byte {
      switch string(match[1]) {
      case "i":
        return []byte(ip)
      case "t":
        return []byte(time.Now().Format("2/Jan/2006:15:04:05 -0700"))
      case "m":
        return []byte(r.Method)
      case "u":
        return []byte(r.URL.String())
      case "s":
        return []byte(strconv.Itoa(traceWriter.statusCode))
      case "b":
        return []byte(formatSize(traceWriter.bytesWritten))
      case "a":
        return []byte(r.Header.Get("User-Agent"))
      case "%":
        return []byte("%")
      default:
        return match
      }
    }))

    fmt.Println(line)
  })
}

type NoCacheResponseWriter struct {
  impl http.ResponseWriter
}

func (writer *NoCacheResponseWriter) Header() http.Header {
  header := writer.impl.Header()
  header.Set("Cache-Control", "private, max-age=0, no-cache")
  header.Del("Last-Modified")
  header.Del("ETag")
  return header
}

func (writer *NoCacheResponseWriter) Write(bytes []byte) (int, error) {
  return writer.impl.Write(bytes)
}

func (writer *NoCacheResponseWriter) WriteHeader(statusCode int) {
  writer.impl.WriteHeader(statusCode)
}

func NoCacheServer(h http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    r.Header.Set("Cache-Control", "no-cache")
    r.Header.Del("If-Modified-Since")
    r.Header.Del("If-None-Match")
    h.ServeHTTP(&NoCacheResponseWriter{w}, r)
  })
}

func main() {
  app := cli.App("sfs", "Static file server - https://github.com/schmich/sfs")
  app.Version("v version", "sfs " + Version)

  port := app.IntOpt("p port", 8080, "Listening port")
  iface := app.StringOpt("i iface interface", "127.0.0.1", "Listening interface")
  allIface := app.BoolOpt("g global", false, "Listen on all interfaces (overrides -i)")
  dir := app.StringOpt("d dir directory", ".", "Directory to serve")
  browser := app.BoolOpt("b browser", false, "Launch web browser")
  trace := app.StringOpt("t trace", "", "Trace format (%i %t %m %u %s %b %a)")
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
      server = NoCacheServer(server)
    }

    if strings.TrimSpace(*trace) != "" {
      server = TraceServer(server, *trace)
    }

    fmt.Printf(">> Serving %s\n", *dir)
    fmt.Printf(">> Listening on %s\n", listen)
    fmt.Println(">> Ctrl+C to stop")

    if *browser {
      url := "http://127.0.0.1" + portPart
      open.Start(url)
    }

    panic(http.ListenAndServe(listen, server))
  }

  app.Run(os.Args)
}
