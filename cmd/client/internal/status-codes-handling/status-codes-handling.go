package status_codes_handling

import (
	"bytes"
	sc "github.com/AlisherFozilov/file-server/pkg/status-codes"
)

func IsStatusCodeError(statusCodeGot []byte) bool {
	if bytes.Equal(statusCodeGot, sc.OK) {
		return false
	}
	return true
}

func HandleStatusCode(statusCodeGot []byte) string {
	switch {
	case cmp(sc.OK, statusCodeGot):
		return "OK"
	case cmp(sc.FILE_NOT_EXISTS, statusCodeGot):
		return "file does not exist on server"
	default:
		return "unknown error code"
	}
}

func cmp(statCode []byte, statusCodeGot []byte) bool {
	return bytes.Equal(statCode, statusCodeGot)
}
