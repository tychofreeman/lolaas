package lolaas

import (
    "fmt"
    "net/http"
    "strings"
    "regexp"
)

var lolRegex *regexp.Regexp
var loRegex *regexp.Regexp
var olthRegex *regexp.Regexp
func init() {
    http.HandleFunc("/lol/", lolHandler)
    http.HandleFunc("/", handler)
    lolRegex, _ = regexp.Compile("(.*[abcdfghkpst])o([^l].*)")
    loRegex, _ = regexp.Compile("(.*[abcdfghkpst])(ol.*)")
    olthRegex, _ = regexp.Compile("(.*bo)(th.*)")
}

func lolify(in string) string {
    out := olthRegex.ReplaceAllString(in, "${1}l${2}")   
    if out == in {
        out = lolRegex.ReplaceAllString(in, "${1}lol${2}")
        if out == in {
            out = loRegex.ReplaceAllString(in, "${1}l${2}")
        }
    }
    return out
}

func lolHandler(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) > 2 {
        out := lolify(parts[2])
        fmt.Fprintf(w, out)
    } else {
        handler(w, r)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, home)
}

var home = `
<html>
<head><title>LOL As A Service</title></head>
<body>
<h2>LOL</h2>
<p>LOL As A Service lets you write like @computionist!</p>
<h2>API</h2>
<h3>/lol/:word</h3>
<p>This will find the best fitting transformation like:</p>
<ul>
<li>both -> bolth</li>
<li>python -> pythloln</li>
<li>dolt -> dlolt</li>
</ul>
</body>
</html>
`
