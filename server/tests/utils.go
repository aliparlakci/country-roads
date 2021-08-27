package tests

import "io"

func BodyReader(reader io.Reader) ([]byte, error) {
	rawBody := make([]byte, 0)
	chunk := make([]byte, 8)
	for {
		n, err := reader.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		rawBody = append(rawBody, chunk[:n]...)
	}

	return rawBody, nil
}
