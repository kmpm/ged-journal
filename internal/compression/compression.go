package compression

import (
	"bytes"
	"compress/zlib"
	"io"
)

// Inflate decompresses a zlib compressed byte slice.
// From small to big
func Inflate(raw []byte) ([]byte, error) {

	r, err := zlib.NewReader(bytes.NewBuffer(raw))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	p := make([]byte, 1024)
	out := new(bytes.Buffer)
	for {
		n, err := r.Read(p)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		out.Write(p[:n])
	}

	return out.Bytes(), nil
}

// Deflate compresses a byte slice using zlib compression.
// From big to small
func Deflate(raw []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(raw)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
