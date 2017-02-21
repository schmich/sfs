# sfs

A go-based static file web server for serving files from a directory on macOS, Windows, or Linux.

Intended as a lightweight development tool for viewing static sites, e.g. documentation, blogs, diagnostic reports, HTML mockups, and prototypes.

## Install

- [Download the zero-install binary](https://github.com/schmich/sfs/releases) to a directory on your `PATH`; or
- `go get -u github.com/schmich/sfs/... && go install github.com/schmich/sfs/...`

## Usage

```
Usage: sfs [-p=<port>] [-i=<interface>] [-s] [-a [USER] PASS] [-g] [-d=<dir>] [-b] [-l=<format>] [-q] [-c]

Static File Server - https://github.com/schmich/sfs

Arguments:
  USER=""      Username for digest authentication
  PASS=""      Password for digest authentication

Options:
  -p, --port=8080                        Listening port
  -i, --iface, --interface="127.0.0.1"   Listening interface
  -s, --secure=false                     Enable HTTPS with self-signed TLS certificate
  -a, --auth=false                       Enable HTTP digest authentication
  -g, --global=false                     Listen on all interfaces (overrides -i)
  -d, --dir, --directory="."             Directory to serve
  -b, --browser=false                    Launch web browser
  -l, --log="%i - %m %u %s"              Log format: %i %t %m %u %s %b %a
  -q, --quiet=false                      Disable request logging
  -c, --cache=false                      Allow cached responses
  -v, --version                          Show the version and exit
```

## Examples

Start a web server for files in the current directory and launch the default browser:

```
sfs -b
```

Specify a port:

```
sfs -p 777
```

Allow external connections:

```
sfs -i 0.0.0.0
sfs -g
```

Serve files from another directory:

```
sfs -d ../bloop
```

## Advanced

### HTTPS with SSL/TLS

Serve content via HTTPS with a self-signed TLS certificate:

```
sfs -s
```

The TLS certificate is randomly generated at startup. Browsers will warn you about an insecure connection since the certificate is self-signed.

### Digest Authentication

Enable HTTP digest authentication with a username and password:

```
sfs -a gordon p4ssw0rd
```

Username is optional. Password is required. If a username is not specified, any non-empty username will work. A password of `-` will prompt you for the password via stdin:

```
sfs -a -
```

### Logging

Log requests with `-l`:

```bash
sfs -l "%i - [%t] %m %u %s %b - %a"
# 127.0.0.1 - [21/Jul/2016:21:07:51 -0500] GET / 200 273 - Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36
```

Log format:

- `%i` Remote IP address
- `%t` Request timestamp
- `%m` Request HTTP method
- `%u` Request URL
- `%s` Response status code
- `%b` Response length (bytes)
- `%a` Request user agent (`User-Agent` HTTP header)

### Caching

By default, `sfs` modifies incoming and outgoing cache headers (`Cache-Control`, `If-None-Match`, `If-Modified-Since`, `Last-Modified`, `ETag`) to ensure no caching occurs. To allow caching, this can be disabled with:

```
sfs -c
```

## License

Copyright &copy; 2016 Chris Schmich  
MIT License. See [LICENSE](LICENSE) for details.
