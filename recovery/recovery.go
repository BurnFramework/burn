package recovery

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
)

const (
	panicHtml = `<html>
<head><title>PANIC: %s</title>
<style type="text/css">
html, body {
	font-family: "Roboto", sans-serif;
	color: #333333;
	background-color: #ea5343;
	margin: 0px;
}
h1 {
	color: #d04526;
	background-color: #ffffff;
	padding: 20px;
	border-bottom: 1px dashed #2b3848;
}
pre {
	margin: 20px;
	padding: 20px;
	border: 2px solid #2b3848;
	background-color: #ffffff;
}
</style>
</head><body>
<h1>PANIC</h1>
<pre style="font-weight: bold;">%s</pre>
<pre>%s</pre>
</body>
</html>`
)

var (
	dunno     = []byte("???")
	centerDot = []byte(".")
	dot       = []byte(".")
	slash     = []byte("/")
)

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, lin, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}
