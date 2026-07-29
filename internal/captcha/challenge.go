package captcha

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"strings"
	"time"
)

var (
	ErrChallenge = errors.New("验证码 challenge 不合法")
	ErrNotFound  = errors.New("验证码 challenge 不存在")
	ErrExpired   = errors.New("验证码 challenge 已过期")
	ErrAnswer    = errors.New("验证码答案不正确")
)

// Challenge 是可一次性消费的验证码状态。
type Challenge struct {
	ID        string    `json:"id"`
	Digest    [32]byte  `json:"digest"`
	ExpiresAt time.Time `json:"expires_at"`
}

// New 创建只保存摘要的验证码 challenge。
func New(id, answer, salt string, expiresAt time.Time) (Challenge, error) {
	if strings.TrimSpace(id) != id || id == "" || len(id) > 128 ||
		answer == "" || len(answer) > 256 || len(salt) < 16 || expiresAt.IsZero() {
		return Challenge{}, ErrChallenge
	}
	return Challenge{
		ID:        id,
		Digest:    digest(answer, salt),
		ExpiresAt: expiresAt,
	}, nil
}

func (p Challenge) verify(answer, salt string, now time.Time) error {
	if p.ID == "" || p.ExpiresAt.IsZero() {
		return ErrChallenge
	}
	if !now.Before(p.ExpiresAt) {
		return ErrExpired
	}
	actual := digest(answer, salt)
	if subtle.ConstantTimeCompare(actual[:], p.Digest[:]) != 1 {
		return ErrAnswer
	}
	return nil
}

func digest(answer, salt string) [32]byte {
	return sha256.Sum256([]byte(salt + "\x00" + answer))
}
