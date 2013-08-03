package lolaas

import (
    "fmt"
    "net/http"
    "strings"
    "regexp"
)

var lolRegex *regexp.Regexp
var olthRegex *regexp.Regexp
func init() {
    http.HandleFunc("/lol/", lolHandler)
    http.HandleFunc("/", handler)
    lolRegex, _ = regexp.Compile("(.*[abcdfghkpst])o(.*)")
    olthRegex, _ = regexp.Compile("(.*bo)(th.*)")
}

func lolify(in string) string {
    out := olthRegex.ReplaceAllString(in, "${1}l${2}")   
    if out == in {
        out = lolRegex.ReplaceAllString(in, "${1}lol${2}")
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
    fmt.Fprintf(w, "LOL!!")
}
