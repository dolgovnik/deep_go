package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	storage []error
}

func (e *MultiError) Error() string {
	if len(e.storage) == 0 {
		return ""
	}
	
	res := fmt.Sprintf("%v errors occured:\n", len(e.storage))
	for _, er := range e.storage {
		res = fmt.Sprintf("%v\t* %v", res,  er.Error())
	}
	res = fmt.Sprintf("%v\n", res)
	return res
}

func (e *MultiError) add(err error) {
	if err != nil {
		var multiEr *MultiError
		if errors.As(err, &multiEr) {
			e.storage = append(e.storage, multiEr.storage...)
		} else {
			e.storage = append(e.storage, err)
		}
	}
}

func Append(err error, errs ...error) *MultiError {
	tempErr := &MultiError{}
	tempErr.add(err)
	for _, e := range errs {
		tempErr.add(e)
	}
	fmt.Printf("%v\n", tempErr.storage)
	return tempErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
