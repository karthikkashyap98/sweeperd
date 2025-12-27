package actions

import (
	"fmt"

	"github.com/karthikkashyap98/sweeperd/internal/rules"
	"github.com/karthikkashyap98/sweeperd/internal/utils"
)

type RuleFunction func(spec rules.Rule, m Matcher) (Action, error)

var registry = map[rules.ActionType]RuleFunction{
	rules.Move: newMove,
}

func Register(kind rules.ActionType, c RuleFunction) {
	registry[kind] = c
}

func NewAction(spec rules.Rule, m Matcher) (Action, error) {
	ctor, ok := registry[spec.Action.Type]
	if !ok {
		return nil, fmt.Errorf("unknown action type: %s", spec.Action.Type)
	}
	return ctor(spec, m)
}

func newMove(spec rules.Rule, m Matcher) (Action, error) {
	dst := utils.ExpandHome(spec.Action.Target)
	src := utils.ExpandHome(spec.Match.Folder)

	return &MoveInstruction{
		Source:      src,
		Destination: dst,
		Matcher:     m,
	}, nil
}
