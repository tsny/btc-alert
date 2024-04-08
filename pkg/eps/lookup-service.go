package eps

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	inner map[*Security]*InfoBall
	mutex *sync.Mutex
}

type InfoBall struct {
	Publisher *Publisher
	Queue     *CandleStack
}

func NewSecurityLookup() *Service {
	return &Service{make(map[*Security]*InfoBall), &sync.Mutex{}}
}

func (s *Service) Register(sec *Security, pub *Publisher, queue *CandleStack) *InfoBall {
	found := s.FindSecurityByNameOrTicker(sec.Name)
	if found != nil {
		return nil
	}
	info := &InfoBall{pub, queue}
	s.mutex.Lock()
	s.inner[sec] = info
	s.mutex.Unlock()
	log.Infof("Tracking security %s [%s] (%s) - [%s]", sec.Name, sec.Source, sec.Ticker, sec.Type.String())
	return info
}

func (s *Service) FindSecurityByNameOrTicker(name string) *InfoBall {
	name = strings.ToLower(name)
	for k, v := range s.inner {
		if name == strings.ToLower(k.Name) || name == strings.ToLower(k.Ticker) {
			return v
		}
		// look thru additional names
		for _, addl := range k.AdditionalNames {
			if name == strings.ToLower(addl) {
				return v
			}
		}
	}
	return nil
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

func (s *Service) GetAllTracked() []*InfoBall {
	var arr []*InfoBall
	for _, v := range s.inner {
		arr = append(arr, v)
	}
	return arr
}
