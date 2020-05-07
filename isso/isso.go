package isso

import (
	"github.com/gorilla/securecookie"
	"wrong.wang/x/go-isso/config"
)

// ISSO do the main logical staff
type ISSO struct {
	storage Storage
	config  config.Config
	guard   guard
}

type guard struct {
	sc *securecookie.SecureCookie
}

// New a ISSO instance
func New(cfg config.Config, storage Storage) *ISSO {
	var HashKey, BlockKey string
	var err error
	if HashKey, err = storage.GetPreference("hask-key"); err != nil {
		HashKey = string(securecookie.GenerateRandomKey(64))
	}
	if BlockKey, err = storage.GetPreference("block-key"); err != nil {
		BlockKey = string(securecookie.GenerateRandomKey(32))
	}
	return &ISSO{
		config: cfg,
		guard: guard{
			sc: securecookie.New([]byte(HashKey), []byte(BlockKey)),
		},
		storage: storage,
	}
}
