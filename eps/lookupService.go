package eps

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	inner map[*Security]*Publisher
	mutex *sync.Mutex
}

func NewSecurityLookup() *Service {
	return &Service{make(map[*Security]*Publisher), &sync.Mutex{}}
}

func (s *Service) Register(sec *Security, pub *Publisher) bool {
	found, _ := s.FindSecurityByNameOrTicker(sec.Name)
	if found != nil {
		return false
	}
	s.mutex.Lock()
	s.inner[sec] = pub
	s.mutex.Unlock()
	log.Infof("Tracking security %s [%s] - [%s]", sec.Name, sec.Source, sec.Type.String())
	return true
}

func (s *Service) FindSecurityByNameOrTicker(name string) (*Security, *Publisher) {
	name = strings.ToLower(name)
	for k, v := range s.inner {
		if name == strings.ToLower(k.Name) || name == strings.ToLower(k.Ticker) {
			return k, v
		}
		// look thru additional names
		for _, addl := range k.AdditionalNames {
			if name == strings.ToLower(addl) {
				return k, v
			}
		}
	}
	return nil, nil
}

func (s *Service) FindSecurityByTicker(ticker string) *Security {
	ticker = strings.ToLower(ticker)
	for k := range s.inner {
		if ticker == k.Ticker {
			return k
		}
	}
	return nil
}
