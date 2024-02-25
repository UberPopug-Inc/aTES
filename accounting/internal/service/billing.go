package service

import (
	"time"
)

type Record struct {
	Timestamp time.Time
	PopugID   string
	debit     uint64
	credit    uint64
}

type BillingStorage struct {
	log []Record
}

func NewPopugsStorage() *BillingStorage {
	return &BillingStorage{log: []Record{}}
}

func (s *BillingStorage) Debit(popugID string, value uint64) {
	s.log = append(s.log, Record{
		Timestamp: time.Now(),
		PopugID:   popugID,
		debit:     value,
	})
}

func (s *BillingStorage) Credit(popugID string, value uint64) {
	s.log = append(s.log, Record{
		Timestamp: time.Now(),
		PopugID:   popugID,
		credit:    value,
	})
}
