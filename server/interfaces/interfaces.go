package interfaces

import "context"

type Validator interface {
	SetDto(dto interface{})
	Validate(ctx context.Context) (bool, error)
}

type Jsonifier interface {
	Jsonify() map[string]interface{}
}
