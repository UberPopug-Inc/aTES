package service

import "math/rand"

type Popug struct {
	ID   string
	name string
}

type PopugsStorage struct {
	popugs map[string]Popug
}

func NewPopugsStorage() *PopugsStorage {
	return &PopugsStorage{popugs: make(map[string]Popug)}
}

func (s *PopugsStorage) Set(popug Popug) {
	s.popugs[popug.ID] = popug
}

func (s *PopugsStorage) GetRandom() Popug {
	n := rand.Intn(len(s.popugs))
	for _, popug := range s.popugs {
		if n--; n == 0 {
			return popug
		}
	}
	return Popug{}
}
