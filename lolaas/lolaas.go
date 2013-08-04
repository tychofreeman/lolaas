package lolaas

import (
    "fmt"
    "net/http"
    "strings"
    "html/template"
    "appengine"
    "encoding/json"
    
)

type ContentTypeWriter func(http.ResponseWriter, appengine.Context, interface{})

var contentTypeWriters map[string]ContentTypeWriter
func init() {
    http.HandleFunc("/lol/", lolHandler)
    http.HandleFunc("/jerk/", jerkHandler)
    http.HandleFunc("/", handler)

    contentTypeWriters = map[string]ContentTypeWriter{
        "application/json":writeJSON,
        "application/xml" :writeXML,
    }
}

func writeJSON(w http.ResponseWriter, c appengine.Context, out interface{}) {
    if marshalled, err := json.Marshal(out); err != nil {
        c.Errorf("Trying to Marshal, but got error %v\n", err)
        fmt.Fprintf(w, "{\"err\": \"Could not write requested data - probably because you're a jerk.\"}")
    } else {
        w.Header().Set("Content-Type","application/json")
        w.Write(marshalled)
    }
}

func writeXML(w http.ResponseWriter, c appengine.Context, out interface{}) {
    w.Header().Set("Content-Type","application/xml")
    fmt.Fprintf(w, "<jerk who=\"you\">%s/jerk/You</jerk>", appengine.DefaultVersionHostname(c))
}

func writePlainText(w http.ResponseWriter, c appengine.Context, out interface{}) {
    fmt.Fprintf(w, "%v", out)
}

func writeWithContentType(w http.ResponseWriter, r *http.Request, out interface{}) {
    accept := r.Header.Get("Accept") 
    var writer ContentTypeWriter = writePlainText
    if ctw, err := contentTypeWriters[accept]; err {
        writer = ctw
    }

    writer(w, appengine.NewContext(r), out)
}

func jerkHandler(w http.ResponseWriter, r *http.Request) {
   writeWithContentType(w, r, Jerk{"You", "jerk"})
}

func lolHandler(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) > 2 {
        out := applyBestLol(parts[2])
        writeWithContentType(w, r, out)
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
<li><h3>--default--</h3><p>Given the default Accept header, you should receive a plain text string.</p></li>
<li><h3>application/json</h3><p>Given the Accept header 'application/json', you should receive a JSON string.</p></li>
<li><h3>application/xml</h3><p>Given the Accept header 'application/xml', you should receive a XML string.</p></li>
</ul>
<br/>
</div>
</body>
</html>
`)
