package lolaas

import (
    "fmt"
    "net/http"
    "strings"
    "regexp"
    "html/template"
    "appengine"
)

type loller func(string)(string,bool)

func lollifier(r *regexp.Regexp, template string) loller {
    var l loller = func(in string) (string, bool) {
        out := r.ReplaceAllString(in, template)
        return out, (out != in)
    }
    return l
}

var regexes []loller
func init() {
    http.HandleFunc("/lol/", lolHandler)
    http.HandleFunc("/", handler)
    regexes = []loller {
        lollifier(regexp.MustCompile("(.*bo)(th.*)"), "${1}l${2}"),
        lollifier(regexp.MustCompile("(.*[abcdfghkpst])o([^l].*)"), "${1}lol${2}"),
        lollifier(regexp.MustCompile("(.*[abcdfghkpst])(ol.*)"), "${1}l${2}"),
        lollifier(regexp.MustCompile("(.*)el+"), "${1}lol"),
        lollifier(regexp.MustCompile("(.*[^l])le"), "${1}lol"),
    }
}

type Lollipop struct {
    Input string
    Output string
}

func (l Lollipop) String() string {
    return l.Output
}

func lolify(in string) Lollipop {
    for _, r := range regexes {
        if out, ok := r(in); ok {
            return Lollipop{in, out}
        }
    }
    return Lollipop{in, in}
}

func lolHandler(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) > 2 {
        out := lolify(parts[2])
        fmt.Fprintf(w, "%s", out)
    } else {
        handler(w, r)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    hostname := appengine.DefaultVersionHostname(c)
    
    home.Execute(w, hostname)
}

var home,_ = template.New("home").Parse(`
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
<li>castle -> castlol</li>
<li>haskell -> hasklol</li>
</ul>
<h2>Example:</h2>
<code>
{{.}}/lol/python
</code>
<p>Results in:</p>
<code>
pythloln
</code>
</body>
</html>
`)
