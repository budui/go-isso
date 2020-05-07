package isso

import (
	"encoding/json"
	"fmt"
	"io"
)

func decodeComment(r io.ReadCloser) (submittedComment, error) {
	defer r.Close()

	var s submittedComment
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&s); err != nil {
		return submittedComment{}, fmt.Errorf("invalid JSON payload: %v", err)
	}

	return s, nil
}
