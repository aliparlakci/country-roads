package common

var MinDate int64 = -62135596800

func NoError(err error, callback func()) {
	if err != nil {
		callback()
	}
}
