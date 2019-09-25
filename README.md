# sfs

An HTTP static file web server and proxy for serving files from a directory on macOS, Windows, or Linux.

Intended as a lightweight development tool for viewing static sites, e.g. documentation, blogs, diagnostic reports, HTML mockups, and prototypes.

## Install

[Download the zero-install binary](https://github.com/schmich/sfs/releases) to a directory on your `PATH`.

## Usage

```
Usage: sfs [-l=<ip:port>] [-s] [-a [USER] PASS] [-d=<dir>] [-b] [-f=<format>] [-q] [-c] [-x=<url>]

Static File Server - https://github.com/schmich/sfs

Arguments:
  USER            Username for digest authentication
  PASS            Password for digest authentication

Options:
  -l, --listen    IP and port to listen on (default "127.0.0.1:8080")
  -s, --secure    Enable HTTPS with self-signed TLS certificate
  -a, --auth      Enable digest authentication
  -d, --dir       Directory to serve
  -b, --browser   Open web browser after server starts
  -f, --format    Log format: %i %t %m %u %s %b %a (default "%i - %m %u %s")
  -q, --quiet     Disable request logging
  -c, --cache     Allow cached responses
  -x, --proxy     Proxy requests to upstream server (implies -c)
  -v, --version   Show the version and exit
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
sfs -l 0.0.0.0
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

### Proxying Requests

`sfs` can act as an HTTP proxy. This is an altogether separate mode of operation from serving static files. This enables you to use `sfs` as a TLS-secured, digest-authenticated, logging frontend for another development server.

```
sfs -x localhost:4567
```

### Logging

Change request logging format with `-f`:

```bash
sfs -f "%i - [%t] %m %u %s %b - %a"
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

Copyright &copy; 2016 Chris Schmich  \
MIT License. See [LICENSE](LICENSE) for details.
