package actions

import (
	"fmt"

	"github.com/karthikkashyap98/sweeperd/internal/rules"
	"github.com/karthikkashyap98/sweeperd/internal/utils"
)

type Constructor func(spec rules.Action, src, dst string, m Matcher) (Action, error)

var registry = map[rules.ActionType]Constructor{
	rules.Move: newMove,
}

func Register(kind rules.ActionType, c Constructor) {
	registry[kind] = c
}

func NewAction(spec rules.Action, source, destination string, m Matcher) (Action, error) {
	ctor, ok := registry[spec.Type]
	if !ok {
		return nil, fmt.Errorf("unknown action type: %s", spec.Type)
	}
	return ctor(spec, source, destination, m)
}

func newMove(spec rules.Action, src, dst string, m Matcher) (Action, error) {
	if dst == "" {
		dst = spec.Target
	}
	return &MoveInstruction{
		Source:      utils.ExpandHome(src),
		Destination: utils.ExpandHome(dst),
		Matcher:     m,
	}, nil
}
