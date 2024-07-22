package testlib

import "fmt"

func Format(sut, expected, actual, input any) string {
	return fmt.Sprintf("\n expected: %v\n actual  : %v\n input   : %#v\n sut     : %#v", expected, actual, input, sut)
}

func FormatWithoutInput(sut, expected, actual any) string {
	return Format(sut, expected, actual, noInput)
}

func FormatError(err error, sut, input any) string {
	return fmt.Sprintf("\n input   : %#v\n sut     : %#v\n error: %+v", input, sut, err)
}

func FormatErrorWithoutInput(err error, sut any) string {
	return FormatError(err, sut, noInput)
}

const noInput = "<N/A>"
