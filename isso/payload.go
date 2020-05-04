package isso

import (
	"encoding/json"
	"fmt"
	"io"

	"wrong.wang/x/go-isso/isso/model"
)

func decodeAcceptComment(r io.ReadCloser) (model.SubmitComment, error) {
	defer r.Close()

	var s model.SubmitComment
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&s); err != nil {
		return model.SubmitComment{}, fmt.Errorf("invalid JSON payload: %v", err)
	}

	return s, nil
}
