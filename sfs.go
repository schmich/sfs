package main

import (
  "os"
  "fmt"
  "time"
  "strconv"
  "strings"
  "regexp"
  "net/http"
  "crypto/tls"
  "path/filepath"
  "github.com/jawher/mow.cli"
  "github.com/skratchdot/open-golang/open"
  "github.com/mh-cbon/gssc"
  "github.com/abbot/go-http-auth"
  "golang.org/x/crypto/ssh/terminal"
)

type LogResponseWriter struct {
  impl http.ResponseWriter
  bytesWritten int
  statusCode int
}

func NewLogResponseWriter(writer http.ResponseWriter) *LogResponseWriter {
  return &LogResponseWriter{
    impl: writer,
    bytesWritten: 0,
    statusCode: 200,
  }
}

func (writer *LogResponseWriter) Header() http.Header {
  return writer.impl.Header()
}

func (writer *LogResponseWriter) Write(bytes []byte) (int, error) {
  writer.bytesWritten += len(bytes)
  return writer.impl.Write(bytes)
}

func (writer *LogResponseWriter) WriteHeader(statusCode int) {
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

func LogHandler(h http.Handler, log string) http.Handler {
  formatter, _ := regexp.Compile("%.")

  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    logWriter := NewLogResponseWriter(w)
    h.ServeHTTP(logWriter, r)

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
        return []byte(strconv.Itoa(logWriter.statusCode))
      case "b":
        return []byte(formatSize(logWriter.bytesWritten))
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

func NoCacheHandler(h http.Handler) http.Handler {
  return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    r.Header.Set("Cache-Control", "no-cache")
    r.Header.Del("If-Modified-Since")
    r.Header.Del("If-None-Match")
    h.ServeHTTP(&NoCacheResponseWriter{w}, r)
  })
}

func AuthHandler(h http.Handler, realm, username, password string) http.Handler {
  handler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
    h.ServeHTTP(w, r)
  })

  authenticator := auth.NewDigestAuthenticator(realm, func (user, realm string) string {
    return password
  })

  authenticator.PlainTextSecrets = true

  return authenticator.Wrap(func (w http.ResponseWriter, r *auth.AuthenticatedRequest) {
    if (username != "") && (r.Username != username) {
      authenticator.RequireAuth(w, &r.Request)
    } else {
      handler(w, &r.Request)
    }
  })
}

func readPassword(prompt string) string {
  fmt.Print(prompt)
  password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
  fmt.Println()
  if err != nil {
    panic(err)
  }
  return string(password)
}

func main() {
  app := cli.App("sfs", "Static File Server - https://github.com/schmich/sfs")
  app.Spec = "[-p=<port>] [-i=<interface>] [-s] [-a [USER] PASS] [-g] [-d=<dir>] [-b] [-l=<format>] [-q] [-c]"

  port := app.IntOpt("p port", 8080, "Listening port")
  iface := app.StringOpt("i iface interface", "127.0.0.1", "Listening interface")
  secure := app.BoolOpt("s secure", false, "Enable HTTPS with self-signed TLS certificate")
  auth := app.BoolOpt("a auth", false, "Enable HTTP digest authentication")
  authUser := app.StringArg("USER", "", "Username for digest authentication")
  authPass := app.StringArg("PASS", "", "Password for digest authentication")
  allIface := app.BoolOpt("g global", false, "Listen on all interfaces (overrides -i)")
  dir := app.StringOpt("d dir directory", ".", "Directory to serve")
  browser := app.BoolOpt("b browser", false, "Launch web browser")
  log := app.StringOpt("l log", "%i - %m %u %s", "Log format: %i %t %m %u %s %b %a")
  quiet := app.BoolOpt("q quiet", false, "Disable request logging")
  cache := app.BoolOpt("c cache", false, "Allow cached responses")

  app.Version("v version", "sfs " + Version)

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

    handler := http.FileServer(http.Dir(*dir))
    if !*cache {
      handler = NoCacheHandler(handler)
    }

    if *quiet {
      *log = ""
    }

    withAuth := ""
    if *auth {
      withAuth = " with digest authentication"
      if *authPass == "-" {
        *authPass = readPassword("HTTP digest authentication password? ")
      }
      handler = AuthHandler(handler, *iface, *authUser, *authPass)
    }

    if strings.TrimSpace(*log) != "" {
      handler = LogHandler(handler, *log)
    }

    protocol := "http"
    if *secure {
      protocol = "https"
    }

    fmt.Printf(">> Serving %s\n", *dir)
    fmt.Printf(">> Listening on %s://%s%s\n", protocol, listen, withAuth)
    fmt.Printf(">> Ctrl+C to stop\n")

    if *browser {
      url := "http://127.0.0.1" + portPart
      open.Start(url)
    }

    server := &http.Server{
      Addr: listen,
      Handler: handler,
    }

    if *secure {
      server.TLSConfig = &tls.Config{
        InsecureSkipVerify: true,
        GetCertificate: gssc.GetCertificate(*iface),
      }
      panic(server.ListenAndServeTLS("", ""))
    } else {
      panic(server.ListenAndServe())
    }
  }

  app.Run(os.Args)
}
