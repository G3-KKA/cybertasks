package logger

import (
	"bytes"
	"io"
	"net/http"
)

const contentPlain = "text/plain; charset=UTF-8"

var _ io.Writer = (*httpWriter)(nil)

type httpWriter struct {
	url string
}

// Write implements io.Writer.
func (h *httpWriter) Write(p []byte) (int, error) {
	r := bytes.NewReader(p)
	_, err := http.Post(h.url, contentPlain, r)
	read := len(p) - r.Len()
	if err != nil {
		return read, err
	}

	return read, nil
}

func newHTTPSingleWriter(url string) (*httpWriter, error) {

	return &httpWriter{
		url: url,
	}, nil
}
