package lolaas

type Jerk struct {
    Who string
    Type string
}

func (j Jerk) String() string {
    return j.Who + " are a " + j.Type
}
