package resolver

import (
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

type Resolver struct {
	table        ResolveTable
	global       *Scope
	currentScope *Scope
	builtins     map[string]int
}

type Scope struct {
	parent   *Scope
	depth    int
	symbols  map[string]*Symbol
	nextSlot int
}

func newScope(parent *Scope) *Scope {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}
	return &Scope{
		parent:  parent,
		depth:   depth,
		symbols: map[string]*Symbol{},
	}
}

type Symbol struct {
	idNodeId parser.IdId
	name     string
	kind     SymbolKind
	slot     int
	scope    *Scope
}

type SymbolKind uint8

const (
	SymbolVar SymbolKind = iota
	SymbolFunc
	SymbolParam
	SymbolBuiltin
)

func NewResolver() *Resolver {
	r := &Resolver{
		table:    ResolveTable{},
		builtins: map[string]int{},
	}
	r.global = newScope(nil)
	r.currentScope = r.global
	r.preludeBuiltins()
	return r
}

func (r *Resolver) preludeBuiltins() {
	builtins := []string{
		"newError",
		"errString",
		"len",
		"scan",
		"print",
		"panic",
	}
	for i, name := range builtins {
		r.builtins[name] = i
		r.global.symbols[name] = &Symbol{
			name:     name,
			kind:     SymbolBuiltin,
			idNodeId: parser.IdId(-1),
			// 빌트인은 전역 슬롯에 위치하지 않음
			// Eval과정에서 빌트인 전용 환경의
			// builtInEnv.Slot (=[r.builtIns[name]])에
			// 해당 빌트인을 삽입할 것임
			// 또한, 빌트인 리졸빙 시에는
			// slot값이 r.builtIns[name]이 되도록 해서
			// 이 코드와 합치되도록 하였음
			// Eval단계에선, 빌트인 만날 시에
			// 전역 스코프가 아닌, 빌트인 스코프에서
			// 해당 리졸브의 slot을 찾아서 적용하면 됨.
			slot:  -1,
			scope: r.global,
		}
	}
}

func (r *Resolver) pushScope() {
	r.currentScope = newScope(r.currentScope)
}

func (r *Resolver) popScope() {
	if r.currentScope.parent != nil {
		r.currentScope = r.currentScope.parent
	}
}

func (r *Resolver) isBuiltinName(name string) bool {
	_, ok := r.builtins[name]
	return ok
}

func (r *Resolver) declare(name string, kind SymbolKind, idnodeId parser.IdId) (*Symbol, error) {
	if r.isBuiltinName(name) {
		return nil, fmt.Errorf("builtin name is reserved: %s", name)
	}

	// 셰도잉 허용함
	// 같은 스코프에선 중복을 금지하지만
	// 스코프 다르다면 셰도잉 허용
	if _, exists := r.currentScope.symbols[name]; exists {
		return nil, fmt.Errorf("duplicate declaration: %s", name)
	}
	slot := r.currentScope.nextSlot
	r.currentScope.nextSlot++
	sym := &Symbol{name: name, kind: kind, idNodeId: idnodeId, slot: slot, scope: r.currentScope}
	r.currentScope.symbols[name] = sym
	return sym, nil
}

func (r *Resolver) resolveID(id parser.Id) (ResolvedRef, error) {
	sym := r.lookup(id.Name)
	if sym == nil {
		return ResolvedRef{}, &ResolveError{IdNode: id, Msg: "undefined identifier"}
	}
	distance := r.currentScope.depth - sym.scope.depth
	ref := ResolvedRef{
		Kind:        RefLocal,
		Distance:    distance,
		Slot:        sym.slot,
		RefIdNodeId: sym.idNodeId,
		Name:        sym.name,
	}
	switch sym.kind {
	case SymbolBuiltin:
		ref.Kind = RefBuiltin
		ref.Distance = 0
		// 빌트인의 경우에는
		// ref.Slot에 전역 스코프에서의 slot이 아닌
		// 빌트인 기준의 slot number를 부여함
		// 실행 시에는 symbolBuiltIn만나면
		// 환경이 아닌 "빌트인 환경"의 "빌트인 기준 슬롯"에 접근하도록 함
		ref.Slot = r.builtins[sym.name]
	case SymbolFunc, SymbolVar:
		if sym.scope == r.global {
			ref.Kind = RefGlobal
			ref.Distance = 0
		}
	}
	return ref, nil
}

type ResolveError struct {
	IdNode parser.Id
	Msg    string
}

func (e *ResolveError) Error() string {
	return fmt.Sprintf("resolve error at %s: %s", e.IdNode.String(), e.Msg)
}

func newResolveErr(idNode parser.Id, msg string) *ResolveError {
	return &ResolveError{
		IdNode: idNode,
		Msg:    msg,
	}
}

func (r *Resolver) lookup(name string) *Symbol {
	for scope := r.currentScope; scope != nil; scope = scope.parent {
		if sym, ok := scope.symbols[name]; ok {
			return sym
		}
	}
	return nil
}
