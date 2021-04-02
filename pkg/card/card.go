package card

import (
	"context"
	"errors"
	"strconv"
	"sync"
)

type Card struct {
	Id       int64  `json:"id"`
	Uid      int64  `json:"uid"`
	Number   string `json:"number"`
	CardType string `json:"card_type"`
}

type Service struct {
	Cards []*Card
	mu    sync.RWMutex
}

func NewService() *Service {
	return &Service{}
}

const (
	RealTypeCard         = "real"
	VirtualTypeCard      = "virtual"
	ErrIncorrectCardType = "incorrect card type"
	ErrPermissionDenied  = "uid not have card permission denied"
)

func (s *Service) IssueCard(uid int64, cardType string, newUser bool, ctx context.Context) error {
	if cardType != RealTypeCard && cardType != VirtualTypeCard {
		return errors.New(ErrIncorrectCardType)
	}

	if !s.ContainsUid(uid) && !newUser {
		return errors.New(ErrPermissionDenied)
	}

	Id := int64(len(s.All(ctx)) + 1)
	number := "000" + strconv.FormatInt(Id, 10)

	s.mu.RLock()
	defer s.mu.RUnlock()
	s.Cards = append(s.Cards, &Card{
		Id:       Id,
		Uid:      uid,
		Number:   number,
		CardType: cardType,
	})

	return nil
}

func (s *Service) All(ctx context.Context) []*Card {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Cards
}

func (s *Service) ContainsUid(uid int64) bool {
	for _, card := range s.Cards {
		if uid == card.Uid {
			return true
		}
	}
	return false
}
