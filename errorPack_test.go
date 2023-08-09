package errorPack

import (
	"testing"
)

type TestErrorData struct {
	Message string
}
type StringableTestErrorData struct {
	Message string
}

func (this *StringableTestErrorData) String() string {
	return this.Message
}

func TestNewErrorPack(t *testing.T) {
	data := TestErrorData{
		Message: "test",
	}

	errorPack := NewErrorPack(nil, data)

	if errorPack == nil {
		t.Error("errorPack is nil")
	}

	if errorPack.Data.Message != data.Message {
		t.Error("errorPack.Data.Message is not equal to data.Message")
	}

	if errorPack.InnerError != nil {
		t.Error("errorPack.InnerError is not nil")
	}

	errorPack.InnerError = NewErrorPack(nil, data)

	if errorPack.InnerError == nil {
		t.Error("errorPack.InnerError is nil")
	}

	if errorPack.InnerError.(*ErrorPack[TestErrorData]).Data.Message != data.Message {
		t.Error("errorPack.InnerError.Data.Message is not equal to data.Message")
	}

	errorPack.Data.Message = "test2"
	if errorPack.InnerError.(*ErrorPack[TestErrorData]).Data.Message == errorPack.Data.Message {
		t.Error("errorPack.InnerError.Data.Message is not equal to data.Message")
	}
}

func TestErrorPack_MemberwiseClone(t *testing.T) {
	errorPack := NewErrorPack(nil, TestErrorData{
		Message: "test",
	})

	errorPack2 := errorPack.MemberwiseClone()

	if &errorPack == &errorPack2 {
		t.Error("&errorPack == &errorPack2")
	}

	if &errorPack.Data.Message != &errorPack2.Data.Message {
		t.Error("&errorPack.Data.Message != &errorPack2.Data.Message")
	}

	errorPack.Data.Message = "test2"
	if errorPack.Data.Message != errorPack2.Data.Message {
		t.Error("errorPack.Data.Message != errorPack2.Data.Message")
	}
}

func TestErrorPack_Error(t *testing.T) {
	errorPack1 := NewErrorPack(nil, TestErrorData{
		Message: "test",
	})

	errStr1 := errorPack1.Error()
	if errStr1 != "ErrorPack[*TestErrorData]" {
		t.Error("errorPack1.Error() != \"ErrorPack[TestErrorData]\"")
	}

	errorPack2 := NewErrorPack(nil, StringableTestErrorData{
		Message: "test",
	})

	errStr2 := errorPack2.Error()
	if errStr2 != "ErrorPack[*StringableTestErrorData]\r\ntest" {
		t.Error("errorPack2.Error() != \"ErrorPack[*StringableTestErrorData]\r\ntest\"")
	}
}

func TestErrorPack_StackTrace(t *testing.T) {
	errorPack := NewErrorPack(nil, TestErrorData{
		Message: "test",
	})

	if errorPack.StackTrace() != "" {
		t.Error("errorPack.StackTrace() != \"\"")
	}
}

func TestErrorPack_Throw(t *testing.T) {
	defer func() {
		var r any
		if r = recover(); r == nil {
			t.Error("recover() == nil")
		}

		ep, ok := r.(*ErrorPack[TestErrorData])
		if !ok {
			t.Error("recover() is not *ErrorPack[TestErrorData]")
		}

		if ep.Data.Message != "test" {
			t.Error("rep.Data.Message != \"test\"")
		}

		stackTrace := ep.StackTrace()
		if stackTrace == "" {
			t.Error("rep.StackTrace() == \"\"")
		}

		println(stackTrace)
	}()

	errorPack := NewErrorPack(nil, TestErrorData{
		Message: "test",
	})

	errorPack.Throw()
}
