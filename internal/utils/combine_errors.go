package utils

import (
	"github.com/hashicorp/go-multierror"
)

func CombineErrors(errors ...error) error {
	var res error
	for _, err := range errors {
		if err != nil {
			res = multierror.Append(res, err)
		}
	}
	return res
}
