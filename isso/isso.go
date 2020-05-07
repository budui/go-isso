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
	v  *Validator
	sc *securecookie.SecureCookie
}

// New a ISSO instance
func New(cfg config.Config, storage Storage) *ISSO {
	HashKey, err := storage.GetPreference("hask-key")
	if err != nil {
		HashKey = string(securecookie.GenerateRandomKey(64))
	}
	BlockKey, err := storage.GetPreference("block-key")
	if err != nil {
		BlockKey = string(securecookie.GenerateRandomKey(32))
	}
	return &ISSO{
		config: cfg,
		guard: guard{
			v:  NewValidator(),
			sc: securecookie.New([]byte(HashKey), []byte(BlockKey)),
		},
		storage: storage,
	}
}
