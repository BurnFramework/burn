package burn

import (
	"log"
	"net/http"
	"net/url"
	"path"
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

func prepareStaticOptions(options []StaticOptions) StaticOptions {
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
	return opt
}

func Static(directory string, staticOpt ...StaticOptions) Handler {
	if !filepath.IsAbs(directory) {
		directory = filepath.Join(Root, directory)
	}
	dir := http.Dir(directory)
	opt := prepareStaticOptions(staticOpt)

	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		if req.Method != "GET" && req.Method != "HEAD" {
			return
		}
		if opt.Exclude != "" && strings.HasPrefix(req.URL.Path, opt.Exclude) {
			return
		}
		file := req.URL.Path

		if opt.Prefix != "" {
			if !strings.HasPrefix(file, opt.Prefix) {
				return
			}
			file = file[len(opt.Prefix):]
			if file != "" && file[0] != '/' {
				return
			}
		}
		f, err := dir.Open(file)
		if err != nil {
			if opt.Fallback != "" {
				file = opt.Fallback
				f, err = dir.Open(opt.Fallback)
			}

			if err != nil {
				return
			}
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return
		}

		if fi.IsDir() {
			if !strings.HasSuffix(req.URL.Path, "/") {
				dest := url.URL{
					Path:     req.URL.Path + "/",
					RawQuery: req.URL.RawQuery,
					Fragment: req.URL.Fragment,
				}
				http.Redirect(res, req, dest.String(), http.StatusFound)
				return
			}

			file = path.Join(file, opt.IndexFile)
			f, err = dir.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				return
			}
		}

		if !opt.SkipLogging {
			log.Println("[Static] Serving " + file)
		}

		if opt.Expires != nil {
			res.Header().Set("Expires", opt.Expires())
		}

		http.ServeContent(res, req, file, fi.ModTime(), f)
	}
}
