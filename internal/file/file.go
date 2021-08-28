package file

import (
	rr "github.com/artofey/media-tel-test"
)

type InputFile interface {
	GetRecords() (*rr.Records, error)
	FileName() string
}
