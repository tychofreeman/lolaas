package lolaas

import (
    "fmt"
    "net/http"
    "strings"
    "regexp"
)

var regex *regexp.Regexp
func init() {
    http.HandleFunc("/lol/", lolHandler)
    http.HandleFunc("/", handler)
    regex, _ = regexp.Compile("(.*[abcdfghkpst])o(.*)")
}

func lolify(in string) string {
    return regex.ReplaceAllString(in, "${1}lol${2}")
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
