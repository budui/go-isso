package util

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// Hasher can hash password with pbkdf2
type Hasher struct {
	salt   []byte
	iter   int
	keyLen int
}

// NewHasher return a new Hasher with conf like this : `pbkdf2:arg1:arg2:arg3`.
// The actual conf for NewHasher is pbkdf2:1000:6:sha1, which means 1000
// iterations, 6 bytes to generate and SHA1 as pseudo-random family used for key
// strengthening. Arguments have to be in that order, but can be reduced to
// pbkdf2:4096 for example to override the iterations only.
// TODO: support all possibility for isso config
func NewHasher(conf, salt string) (Hasher, error) {
	r := strings.Split(conf, ":")
	if r[0] != "pbkdf2" {
		return Hasher{}, errors.New("go-isso just support pbkdf2 for now")
	}
	iter := 1000
	keyLen := 6
	var err error

	if len(r) >= 2 {
		iter, err = strconv.Atoi(r[1])
		if err != nil {
			return Hasher{}, fmt.Errorf("arg2:%s", err)
		}
	}
	if len(r) >= 3 {
		keyLen, err = strconv.Atoi(r[1])
		if err != nil {
			return Hasher{}, fmt.Errorf("arg3:%s", err)
		}
	}

	return Hasher{
		iter:   iter,
		keyLen: keyLen,
		salt:   []byte(salt),
	}, nil
}

// Hash hash a string to a hex string.
func (h *Hasher) Hash(p string) string {
	return hex.EncodeToString(pbkdf2.Key([]byte(p), h.salt, h.iter, h.keyLen, sha1.New))
}
