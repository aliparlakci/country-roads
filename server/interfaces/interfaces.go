package interfaces

type Validator interface {
	Validate() (bool, error)
}

type Jsonifier interface {
	Jsonify() map[string]interface{}
}
