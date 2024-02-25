package main

import (
	"github.com/UberPopug-Inc/aTES/accounting/internal/events"
	"github.com/UberPopug-Inc/aTES/accounting/internal/service"
)

func main() {
	k := events.NewKafka()
	s := service.New(k)
	if err := s.Run(); err != nil {
		panic("s.Run()")
	}
}
