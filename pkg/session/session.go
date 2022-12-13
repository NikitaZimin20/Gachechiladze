package session

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
)

type Session struct {
	ID        string
	UserID    uint32
	UserName  string
	UserType  string
	Purchases map[uint32]struct{}
}

func (s *Session) IsPurchased(id uint32) bool {
	_, ok := s.Purchases[id]
	return ok
}

func (s *Session) AddPurchase(id uint32) {
	s.Purchases[id] = struct{}{}
}

func (s *Session) DeletePurchase(id uint32) {
	delete(s.Purchases, id)
}

func NewSession(userID uint32, userName, userType string) *Session {
	// лучше генерировать из заданного алфавита, но так писать меньше и для учебного примера ОК
	randID := make([]byte, 16)
	rand.Read(randID)

	return &Session{
		ID:        fmt.Sprintf("%x", randID),
		UserID:    userID,
		UserName:  userName,
		UserType:  userType,
		Purchases: map[uint32]struct{}{},
	}
}

var (
	ErrNoAuth = errors.New("no session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return &Session{}, ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
