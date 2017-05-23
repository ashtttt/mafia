package topt

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"hash"
	"math"
	"time"
)

// TOPT contains the different configurable values for a given TOTP
type TOPT struct {
	Time     func() time.Time
	TimeStep time.Duration
	Digits   uint8
	Hash     func() hash.Hash
}

// NewTOPT returns an new TOPT object with defualt options
func NewTOPT() *TOPT {
	return &TOPT{
		Time:     time.Now,
		TimeStep: 30 * time.Second,
		Digits:   6,
		Hash:     sha1.New,
	}
}

// DefaultOptions returns an new TOPT object with defualt options
var DefaultOptions = NewTOPT()

// TokenCode generates TOPT code.
func (o *TOPT) TokenCode(secretKey []byte) string {
	if o == nil {
		o = DefaultOptions
	}

	t := o.Time().Unix() / int64(o.TimeStep/time.Second)
	tbuf := make([]byte, 8)

	for i := 7; i >= 0; i-- {
		tbuf[i] = byte(t & 0xff)
		t = t >> 8
	}

	var hashbuf []byte
	hm := hmac.New(o.Hash, secretKey)
	hm.Write([]byte(tbuf))
	hashbuf = hm.Sum(nil)

	offset := int(hashbuf[len(hashbuf)-1] & 0xf)

	code := ((int(hashbuf[offset]) & 0x7f) << 24) |
		((int(hashbuf[offset+1] & 0xff)) << 16) |
		((int(hashbuf[offset+2] & 0xff)) << 8) |
		(int(hashbuf[offset+3]) & 0xff)

	otp := int64(code) % int64(math.Pow10(int(o.Digits)))

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", o.Digits), otp)
}
