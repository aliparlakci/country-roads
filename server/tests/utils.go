package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

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

func IsBodyEqual(expected gin.H, actual io.ReadCloser) (bool, error) {
	actualBytes, err := BodyReader(actual)
	if err != nil {
		return false, err
	}

	parsedExpected := expected

	var parsedActual gin.H
	err = json.Unmarshal(actualBytes, &parsedActual)
	if err != nil {
		return false, nil
	}

	marshalledExpected, _ := json.Marshal(parsedExpected)
	marshalledActual, _ := json.Marshal(parsedActual)

	return string(marshalledExpected) == string(marshalledActual), nil
}
