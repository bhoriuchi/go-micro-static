package static

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"
)

var debugging = false

type static struct {
	dir      string
	fs       http.Handler
	services []string
}

func (s *static) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:   "static_service",
			Usage:  "Service name to exclude from static route matching",
			EnvVar: "STATIC_SERVICE",
		},
		cli.StringFlag{
			Name:   "static_dir",
			Usage:  "Directory to serve static from",
			EnvVar: "STATIC_DIR",
		},
		cli.BoolFlag{
			Name:        "static_debug",
			Usage:       "Log debug output",
			EnvVar:      "STATIC_DEBUG",
			Destination: &debugging,
		},
	}
}

// simple debugger function
func debug(fmtstr string, args ...interface{}) {
	if debugging {
		log.Printf(fmtstr, args...)
	}
}

func (s *static) Commands() []cli.Command {
	return nil
}

func (s *static) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, service := range s.services {
				prefix := fmt.Sprintf("/%s", strings.Trim(service, "/"))
				uri := fmt.Sprintf("/%s", strings.Trim(r.RequestURI, "/"))
				if strings.HasPrefix(uri, prefix) {
					debug("Handled %s with micro handler", r.RequestURI)
					h.ServeHTTP(w, r)
					return
				}
			}

			if _, err := os.Stat(s.dir + r.RequestURI); os.IsNotExist(err) {
				debug("Handled %s with strip prefix handler", r.RequestURI)
				http.StripPrefix(r.RequestURI, s.fs).ServeHTTP(w, r)
				return
			}

			debug("Handled %s with file server", r.RequestURI)
			s.fs.ServeHTTP(w, r)
		})
	}
}

func (s *static) Init(ctx *cli.Context) error {
	s.services = ctx.StringSlice("static_service")
	debug("Ignoring static content for requests on services %v", s.services)
	dir := ctx.String("static_dir")
	if len(dir) == 0 {
		dir = "html"
	}
	s.dir = dir
	s.fs = http.FileServer(http.Dir(s.dir))
	return nil
}

func (s *static) String() string {
	return "static"
}

// NewPlugin returns new plugin
func NewPlugin() plugin.Plugin {
	return &static{}
}
