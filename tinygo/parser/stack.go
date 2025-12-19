package parser

import "errors"

type Stack[V any] struct {
	stack []V
}

func Push[V any](s Stack[V], el V) {
	s.stack = append(s.stack, el)
}

func Pop[V any](s Stack[V]) (V, error) {
	if len(s.stack) == 0 {
		var v V
		return v, errors.New("can't pop empty stack")
	}
	s.stack = s.stack[0 : len(s.stack)-1]
	top := s.stack[len(s.stack)-1]
	return top, nil
}
func GetRollBack[V any](s Stack[V]) func(s Stack[V]) {
	currentTop := len(s.stack) - 1
	return func(s Stack[V]) {
		if currentTop < 0 {
			s.stack = []V{}
			return
		}
		s.stack = s.stack[0 : currentTop+1]
	}
}
