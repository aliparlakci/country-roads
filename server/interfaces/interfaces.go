package interfaces

type Validator interface {
	Validate() (bool, error)
}

type JSONifier interface {
	JSONify() map[string]interface{}
}
