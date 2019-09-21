package auth

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var tokenByUser map[string]Data

func init() {
	tokenByUser = make(map[string]Data)
}

// Data -
type Data struct {
	ID     string
	Token  string
	UserID string
}

func hashToken(account, password string) string {
	s := md5.New()
	s.Write([]byte(fmt.Sprintf("%v+%v + %v", account, password, time.Now().Unix())))
	return hex.EncodeToString(s.Sum(nil))
}

// Login -
func Login(account, password string) (string, error) {
	token := hashToken(account, password)

	tokenByUser[token] = Data{
		ID:     time.Now().Format(time.RFC3339),
		Token:  token,
		UserID: account,
	}

	return token, nil

}

// Logout -
func Logout(token string) error {

	if _, ok := tokenByUser[token]; !ok {
		return errors.New("token not found")
	}
	delete(tokenByUser, token)

	return nil

}

// List -
func List() (int32, map[string]Data) {
	total := int32(len(tokenByUser))

	return total, tokenByUser
}
