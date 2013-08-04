package lolaas

import (
    "regexp"
)


type loller func(string)(string,bool)

func lollifier(r *regexp.Regexp, template string) loller {
    var l loller = func(in string) (string, bool) {
        out := r.ReplaceAllString(in, template)
        return out, (out != in)
    }
    return l
}

var regexes []loller = []loller {
    lollifier(regexp.MustCompile("(.*bo)(th.*)"), "${1}l${2}"),
    lollifier(regexp.MustCompile("(.*[abcdfghkpst])o([^l].*)"), "${1}lol${2}"),
    lollifier(regexp.MustCompile("(.*[abcdfghkpst])(ol.*)"), "${1}l${2}"),
    lollifier(regexp.MustCompile("(.*)el+"), "${1}lol"),
    lollifier(regexp.MustCompile("(.*[^l])le"), "${1}lol"),
}

type Lollipop struct {
    Input string
    Output string
}

func (l Lollipop)String() string{
    return l.Output
}

func applyBestLol(in string) Lollipop {
    for _, r := range regexes {
        if out, ok := r(in); ok {
            return Lollipop{in, out}
        }
    }
    return Lollipop{in, in}
}
