package strutils

import (
	"errors"
)

func IndexOfFirstNonWhitespace(data *[]byte, startIdx uint64) (uint64, error) {
	slicedData := (*data)[startIdx:]
	for i, b := range slicedData {
		if b != ' ' && b != '\t' && b != '\n' && b != '\r' {
			return startIdx + uint64(i), nil
		}
	}
	return 0, errors.New("Cannot find non whitespace character")
}

func IndexOfFirstWhitespace(data *[]byte, startIdx uint64) (uint64, error) {
	slicedData := (*data)[startIdx:]
	for i, b := range slicedData {
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			return startIdx + uint64(i), nil
		}
	}
	return 0, errors.New("Cannot find whitespace character")
}
