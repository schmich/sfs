# sfs

A go-based static file web server for serving files from a directory on macOS, Windows, or Linux.

Intended as a lightweight development tool for viewing static sites, e.g. documentation, blogs, diagnostic reports, HTML mockups, and early prototypes.

## Install

- [Download the zero-install binary](https://github.com/schmich/sfs/releases) to a directory on your `PATH`; or
- `go get -u github.com/schmich/sfs/... && go install github.com/schmich/sfs/...`

## Usage

```
Usage: sfs [OPTIONS]

Static file server - https://github.com/schmich/sfs

Options:
  -v, --version                          Show the version and exit
  -p, --port=8080                        Listening port
  -i, --iface, --interface="127.0.0.1"   Listening interface
  -g, --global=false                     Listen on all interfaces (overrides -i)
  -d, --dir, --directory="."             Directory to serve
  -B, --no-browser=false                 Do not launch browser
  -c, --cache=false                      Allow cached responses
```

Start a web server for files in the current directory:

```
sfs
```

Specify a port:

```
sfs -p 777
```

Allow external connections:

```
sfs -i 0.0.0.0
```

Alternatively:

```
sfs -g
```

Serve files from another directory:

```
sfs -d ../bloop
```

By default, `sfs` modifies incoming and outgoing cache headers (`Cache-Control`, `If-None-Match`, `If-Modified-Since`) to ensure no caching occurs. This can be disabled with:

```
sfs -c
```

## License

Copyright &copy; 2016 Chris Schmich
<br />
MIT License. See [LICENSE](LICENSE) for details.
