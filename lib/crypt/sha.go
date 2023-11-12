package crypt

import (
	"crypto/sha256"
	"encoding/base64"
)

type Sha256Encoder struct{}

func (Sha256Encoder) Encode(mess, key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(mess + key))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

type Enrypter interface {
	Encode(mess, key string) string
}
type Option func(*EncryptOpt)

type EncryptOpt struct {
	key string
}

func WithKey(key string) Option {
	return func(eo *EncryptOpt) {
		eo.key = key
	}
}

func EncryptMess(mess string, encrType Enrypter, opts ...Option) string {
	var option EncryptOpt
	for _, opt := range opts {
		opt(&option)
	}

	encodedMess := encrType.Encode(mess, option.key)
	return encodedMess
}
