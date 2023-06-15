package internal

import (
	"fmt"
	"strings"
)

// InferExt infers the extension from an S3 object key
func InferExt(object_key string) (string, error) {
	ok_splits := strings.Split(object_key, "/")
	object_name := ok_splits[len(ok_splits)-1]
	on_splits := strings.Split(object_name, ".")
	if len(on_splits) < 2 {
		return "", fmt.Errorf("could not figure out the file extension from the object key %s", object_key)
	}
	return on_splits[len(on_splits)-1], nil
}
