package isso

import (
	"encoding/json"
	"fmt"
	"io"

	"wrong.wang/x/go-isso/isso/model"
)

func decodeAcceptComment(r io.ReadCloser) (model.AcceptComment, error) {
	defer r.Close()

	var s model.AcceptComment
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&s); err != nil {
		return model.AcceptComment{}, fmt.Errorf("invalid JSON payload: %v", err)
	}

	return s, nil
}
