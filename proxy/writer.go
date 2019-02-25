package proxy

import (
	"io"
	"net/http"
)

type SequentialWriter struct {
	io.WriterAt
	delegate http.ResponseWriter
}

func NewSequentialWriter(writer http.ResponseWriter) *SequentialWriter {
	return &SequentialWriter{
		delegate: writer,
	}
}

func (writer SequentialWriter) WriteAt(p []byte, off int64) (n int, err error) {
	return writer.delegate.Write(p)
}
