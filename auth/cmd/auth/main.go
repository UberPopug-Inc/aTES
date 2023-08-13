package main

import (
	"github.com/UberPopug-Inc/aTES/auth/internal/service"
)

func main() {
	s := service.New()
	if err := s.Run(); err != nil {
		panic("s.Run()")
	}
}
