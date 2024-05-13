package lib

import "io"

type Decoder interface {
	Init(reader io.Reader) error
	Decode() (*CandidateNode, error)
}
