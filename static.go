package static

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"
)

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
	}
}

func (s *static) Commands() []cli.Command {
	return nil
}

func (s *static) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := os.Stat(s.dir + r.RequestURI); os.IsNotExist(err) {
				for _, service := range s.services {
					prefix := fmt.Sprintf("/%s", strings.Trim(service, "/"))
					if strings.HasPrefix(r.RequestURI, prefix) {
						h.ServeHTTP(w, r)
						return
					}
				}

				http.StripPrefix(r.RequestURI, s.fs).ServeHTTP(w, r)
			} else {
				s.fs.ServeHTTP(w, r)
			}
		})
	}
}

func (s *static) Init(ctx *cli.Context) error {
	s.services = ctx.StringSlice("static_service")
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
