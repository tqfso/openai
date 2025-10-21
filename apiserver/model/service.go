package model

import (
	"apiserver/client/openserver"
)

// 某个模型的所有服务

type Service struct {
	ID          string
	Power       uint64
	Load        uint64
	SelectIndex int
	Targets     []*openserver.ServiceTarget
}

type Services struct {
	SelectIndex int
	Services    []*Service
}

func (s *Service) SelectTarget() *openserver.ServiceTarget {
	s.SelectIndex++
	count := len(s.Targets)
	if count == 0 {
		return nil
	}

	return s.Targets[s.SelectIndex%count]
}

func (s *Services) SelectTarget() *openserver.ServiceTarget {
	s.SelectIndex++
	count := len(s.Services)
	if count == 0 {
		return nil
	}

	service := s.Services[s.SelectIndex%count]

	return service.SelectTarget()
}
