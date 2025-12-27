package resolver

import (
	"fmt"
	"sort"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

type ResolveTable map[parser.IdId]ResolvedRef

type ResolvedRef struct {
	// 참조한 선언 id의 id
	RefIdNodeId parser.IdId
	Kind        RefKind
	//참조한 선언과의 거리
	Distance int
	Slot     int

	//참조한 id의 이름이자, 리졸빙 대상 id의 이름
	Name string
}

type RefKind uint8

const (
	RefLocal RefKind = iota
	RefGlobal
	RefBuiltin
)

func (k RefKind) String() string {
	switch k {
	case RefLocal:
		return "Local"
	case RefGlobal:
		return "Global"
	case RefBuiltin:
		return "Builtin"
	default:
		return "Unknown"
	}
}

func (rt ResolveTable) Print() string {
	if len(rt) == 0 {
		return "<empty>"
	}
	keys := make([]int, 0, len(rt))
	for id := range rt {
		keys = append(keys, int(id))
	}
	sort.Ints(keys)

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		id := parser.IdId(key)
		ref := rt[id]
		line := fmt.Sprintf("#%d %s => #%d kind=%s distance=%d slot=%d", id, ref.Name, ref.RefIdNodeId, ref.Kind.String(), ref.Distance, ref.Slot)
		lines = append(lines, line)
	}
	return parser.JoinLines(lines)
}
