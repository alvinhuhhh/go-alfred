package util

import (
	"crypto/hkdf"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func int64ToBE(v int64) []byte {
	var b [8]byte
	binary.BigEndian.AppendUint64(b[:], uint64(v))
	return b[:]
}

func DeriveDEK(masterKey []byte, keyVersion uint64, chatId int64) ([]byte, error) {
	// generate salt from chatId
	var salt [32]byte
	{
		h := sha256.Sum256(int64ToBE(chatId))
		copy(salt[:], h[:])
	}

	info := []byte(fmt.Sprintf("Alfred-DEK|AES-256-GCM|v=%d|scope=telegram-chat", keyVersion))
	r, err := hkdf.Key(sha256.New, masterKey, salt[:], string(info), 32)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func GetKeyVersion() (uint64, error) {
	value, found := os.LookupEnv("MASTER_KEY_VERSION")
	if !found {
		return 0, errors.New("MASTER_KEY_VERSION is not found")
	}
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetMasterKey(keyVersion uint64) ([]byte, error) {
	name := "MASTER_KEY_V" + strconv.FormatUint(keyVersion, 10)
	encoded, found := os.LookupEnv(name)
	if !found {
		return nil, errors.New(fmt.Sprintf("%s is not found", name))
	}
	value, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return []byte(value), nil
}
