package static

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type StaticOptions struct {
	Prefix      string
	SkipLogging bool
	IndexFile   string
	Expires     func() string
	Fallback    string
	Exclude     string
}

func prepareStaticOptions(options []StaticOptions) {
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}

	if opt.Prefix != "" {
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}
}

/* serves the static */
func Static(directory string, staticOpt ...StaticOptions) Handler {
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(Root, directory)
	}
	dir := http.Dir(directory)
	opt := prepareStaticOptions(staticOpt)
	
	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {

	}

}
