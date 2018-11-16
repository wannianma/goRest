package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/orcaman/concurrent-map"
	"sync"
	"time"
)

var instance *TokenManager
var once sync.Once

type TimeStamp = int64

type innerType struct {
	refreshToken  string
	unixTimeStamp TimeStamp
}

type TokenManager struct {
	cmap.ConcurrentMap
	maxRefreshLife     TimeStamp
	maxAccessTokenLife TimeStamp
	privateKey         *rsa.PrivateKey
	publicKey          *rsa.PublicKey
}

func New(pri_key string, pub_key string, maxAccessTokenSec uint, maxRefreshTokenSec uint) (*TokenManager, error) {
	t := &TokenManager{
		cmap.New(),
		TimeStamp(maxRefreshTokenSec),
		TimeStamp(maxAccessTokenSec),
		nil,
		nil,
	}
	var err error
	t.privateKey, err = parseRsaPrivateKeyFromPemStr(pri_key)
	if err != nil {
		return t, err
	}
	t.publicKey, err = parseRsaPublicKeyFromPemStr(pub_key)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (m *TokenManager) Hit(uid uint, token string) bool {
	var id = fmt.Sprint(uid)
	result := false
	if tmp, ok := m.Get(id); ok {
		inner := tmp.(innerType)
		nowStamp := time.Now().Unix()

		if token == inner.refreshToken && nowStamp < inner.unixTimeStamp {
			// refresh the token live time
			inner.unixTimeStamp = nowStamp + m.maxRefreshLife
			m.Set(id, inner)
			result = true
		}
	}
	return result
}

func (m *TokenManager) Create(uid uint, token string) {
	var id = fmt.Sprint(uid)
	inner := innerType{
		refreshToken:  token,
		unixTimeStamp: time.Now().Unix() + m.maxRefreshLife,
	}
	m.Set(id, inner)
}

func (m *TokenManager) Delete(uid uint) {
	var id = fmt.Sprint(uid)
	m.Remove(id)
}
