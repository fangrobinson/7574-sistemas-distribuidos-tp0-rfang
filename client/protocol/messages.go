package protocol

import (
	"fmt"
)

// GetMessageLength returns the length of the data associated with a specific message code.
func GetMessageLength(code byte) (int, error) {
	switch code {
	case byte(2):
		return 0, nil
	case byte(4):
		return 0, nil
	default:
		return 0, fmt.Errorf("unexpected message code: %v", code)
	}
}
