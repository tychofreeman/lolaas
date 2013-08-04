package lolaas

import (
    "fmt"
    "net/http"
    "strings"
    "regexp"
    "html/template"
    "appengine"
    "encoding/json"
    
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
    http.HandleFunc("/jerk/", jerkHandler)
    http.HandleFunc("/", handler)
    regexes = []loller {
        lollifier(regexp.MustCompile("(.*bo)(th.*)"), "${1}l${2}"),
        lollifier(regexp.MustCompile("(.*[abcdfghkpst])o([^l].*)"), "${1}lol${2}"),
        lollifier(regexp.MustCompile("(.*[abcdfghkpst])(ol.*)"), "${1}l${2}"),
        lollifier(regexp.MustCompile("(.*)el+"), "${1}lol"),
        lollifier(regexp.MustCompile("(.*[^l])le"), "${1}lol"),
    }
}

func jerkHandler(w http.ResponseWriter, r *http.Request) {
   fmt.Fprintf(w, "You are a jerk.")
}

type Lollipop struct {
    Input string
    Output string
}

func lolify(in string) Lollipop {
    for _, r := range regexes {
        if out, ok := r(in); ok {
            return Lollipop{in, out}
        }
    }
    return Lollipop{in, in}
}

func lolJson(w http.ResponseWriter, c appengine.Context, out Lollipop) {
    if marshalled, err := json.Marshal(out); err != nil {
        c.Errorf("Trying to Marshal, but got error %v\n", err)
        fmt.Fprintf(w, "{\"err\": \"Could not write requested data - probably because you're a jerk.\"}")
    } else {
        w.Header().Set("Content-Type","application/json")
        w.Write(marshalled)
    }
}

func lolXml(w http.ResponseWriter, c appengine.Context, out Lollipop) {
    w.Header().Set("Content-Type","application/xml")
    fmt.Fprintf(w, "<jerk who=\"you\">%s/jerk/You</jerk>", appengine.DefaultVersionHostname(c))
}

func lolHandler(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) > 2 {
        out := lolify(parts[2])
        c := appengine.NewContext(r)
        accept := r.Header.Get("Accept") 
        if accept == "application/json" {
            lolJson(w, c, out)
        } else if accept == "application/xml" {
            lolXml(w, c, out)
        } else {
            fmt.Fprintf(w, "%s", out.Output)
        }
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
<!DOCTYPE html>
<html lang="en">
<head>
<title>LOL As A Service</title>
<!--link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.0.0-rc1/css/bootstrap.min.css" -->
<style>
.main{
    top: 5%;
    bottom: 0;
    left: 0;
    right: 0;
    margin:0px auto; 
    width: 50%;
    position: absolute;
    overflow: auto;
    border: solid 1px;
    border-color: #eeeeee;
}

.header{
    top: 0%;
    bottom: 0;
    left: 0;
    right: 0;
    height: 12%;
    background: red;
    margin:0px auto; 
    position: relative;
    overflow: auto;
    border: solid 2px;
    border-color: black;
    text-align: center;
    margin-top: auto;
    margin-bottom: auto;
    display:table;
    width:99%;
}

.big-title {
    vertical-align: middle;
    display:table-cell;
}

</style>
</head>
<body class="container">
<div class="container main">
<div class="header">
<h1 class="big-title">Laugh-Out-Loud As A Service</h1>
</div>
<h2>LOL</h2>
<p>LOL As A Service lets you write like <a href="http://twitter.com/computionist">Doctor Gonzo</a>!</p>
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
<h3>/jerk</h3>
<p>This will print out a string which declares the user to be a jerk.</p>
<h2>Examples:</h2>
<code>
{{.}}/lol/python
</code>
<p>Results in:</p>
<code>
pythloln
</code>
<h2>Supported Accept Headers (Currently working for /lol/ only.)</h2>
<ul>
<li><h3>application/json</h3><p>Given the Accept header 'application/json', you should receive a JSON string.</p></li>
<li><h3>application/xml</h3><p>Given the Accept header 'application/xml', you should receive a XML string.</p></li>
</ul>
</div>
</body>
</html>
`)
