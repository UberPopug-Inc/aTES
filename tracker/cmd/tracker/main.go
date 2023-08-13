package main

import (
	"github.com/UberPopug-Inc/aTES/auth/internal/events"
	"github.com/UberPopug-Inc/aTES/auth/internal/service"
)

func main() {
	k := events.NewKafka()
	s := service.New(k)
	if err := s.Run(); err != nil {
		panic("s.Run()")
	}
}
