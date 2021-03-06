package imaging

import (
	"bytes"
	"image"
	"io"
)

type Processor interface {
	Thumbnail(src io.ReadSeeker) (image.Image, error)
	Fit(src io.ReadSeeker, width, height int) (image.Image, error)
	FitAnimationGIF(src io.Reader, width, height int) (*bytes.Reader, error)
}
