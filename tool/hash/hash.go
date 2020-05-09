package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
	"wrong.wang/x/go-isso/logger"
)

// Worker can hash password with pbkdf2
type Worker struct {
	salt   []byte
	iter   int
	keyLen int
}

// New return a new HashWorker with conf like this : `pbkdf2:arg1:arg2:arg3`.
// The actual conf for NewHashWorker is pbkdf2:1000:6:sha1, which means 1000
// iterations, 6 bytes to generate and SHA1 as pseudo-random family used for key
// strengthening. Arguments have to be in that order, but can be reduced to
// pbkdf2:4096 for example to override the iterations only.
// TODO: support all possibility for isso config
func New(conf, salt string) *Worker {
	r := strings.Split(conf, ":")
	if r[0] != "pbkdf2" {
		logger.Fatal("go-isso just support pbkdf2 for now")
	}
	iter := 1000
	keyLen := 6
	var err error

	if len(r) >= 2 {
		iter, err = strconv.Atoi(r[1])
		if err != nil {
			logger.Fatal("hash conf error - convert arg1(%s) failed: %w", r[1], err)
		}
	}
	if len(r) >= 3 {
		keyLen, err = strconv.Atoi(r[2])
		if err != nil {
			logger.Fatal("hash conf error - convert arg2(%s) failed: %w", r[2], err)
		}
	}

	return &Worker{
		iter:   iter,
		keyLen: keyLen,
		salt:   []byte(salt),
	}
}

// Hash hash a string to a hex string.
func (h *Worker) Hash(p string) string {
	return hex.EncodeToString(pbkdf2.Key([]byte(p), h.salt, h.iter, h.keyLen, sha1.New))
}
