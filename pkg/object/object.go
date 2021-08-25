package object

type Object struct {
	ObjectID   int  `json:"id"`
	Online     bool `json:"online"`
	LastSeen   int64
	ValidUntil int64
}
