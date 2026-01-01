package resolver

import (
	"fmt"
	"sort"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

type HoistInfo struct {
	globalsByName map[string]*Symbol
	globalsById   map[parser.IdId]*Symbol
	// 선언 순서대로 수집
	// var간 초기화 순서는 추후 initOrder에서 보장함
	varOrder     []parser.IdId
	funcOrder    []parser.IdId
	varDeclById  map[parser.IdId]*parser.VarDecl
	funcDeclById map[parser.IdId]*parser.FuncDecl
}

func newHoistInfo() *HoistInfo {
	return &HoistInfo{
		globalsByName: map[string]*Symbol{},
		globalsById:   map[parser.IdId]*Symbol{},
		varDeclById:   map[parser.IdId]*parser.VarDecl{},
		funcDeclById:  map[parser.IdId]*parser.FuncDecl{},
	}
}

func (h *HoistInfo) getById(id parser.IdId) *Symbol {
	return h.globalsById[id]
}

func (h *HoistInfo) getVarDeclById(id parser.IdId) *parser.VarDecl {
	return h.varDeclById[id]
}

func (h *HoistInfo) getFuncDeclById(id parser.IdId) *parser.FuncDecl {
	return h.funcDeclById[id]
}

func (h *HoistInfo) varIds() []parser.IdId {
	return h.varOrder
}

func (h *HoistInfo) funcIds() []parser.IdId {
	return h.funcOrder
}

func (h *HoistInfo) GetVarDeclById(id parser.IdId) *parser.VarDecl {
	return h.getVarDeclById(id)
}

func (h *HoistInfo) GetFuncDeclById(id parser.IdId) *parser.FuncDecl {
	return h.getFuncDeclById(id)
}

func (h *HoistInfo) VarIds() []parser.IdId {
	return h.varIds()
}

func (h *HoistInfo) FuncIds() []parser.IdId {
	return h.funcIds()
}

func (h *HoistInfo) Print() string {
	if h == nil {
		return "<nil hoist>"
	}
	lines := []string{"HoistInfo:"}

	varIds := make([]int, 0, len(h.varOrder))
	for _, id := range h.varOrder {
		varIds = append(varIds, int(id))
	}
	sort.Ints(varIds)
	lines = append(lines, "vars:")
	for _, idv := range varIds {
		id := parser.IdId(idv)
		name := "<missing>"
		if sym := h.getById(id); sym != nil {
			name = sym.name
		}
		lines = append(lines, fmt.Sprintf("  #%d %s", id, name))
	}

	funcIds := make([]int, 0, len(h.funcOrder))
	for _, id := range h.funcOrder {
		funcIds = append(funcIds, int(id))
	}
	sort.Ints(funcIds)
	lines = append(lines, "funcs:")
	for _, idv := range funcIds {
		id := parser.IdId(idv)
		name := "<missing>"
		if sym := h.getById(id); sym != nil {
			name = sym.name
		}
		lines = append(lines, fmt.Sprintf("  #%d %s", id, name))
	}

	return parser.JoinLines(lines)
}

func (r *Resolver) collectPackageDecls(pkg *parser.PackageAST) (*HoistInfo, error) {
	hoist := newHoistInfo()
	for _, decl := range pkg.DeclsOrNil {
		switch node := decl.(type) {
		// 패키지 레벨의 선언들을 돌면서
		// 좌변의 선언 인자들만 글로벌 스코프에 전부 등록
		case *parser.VarDecl:
			for _, id := range node.Ids {
				sym, err := r.declare(id.Name, SymbolVar, id.IdId)
				if err != nil {
					return nil, newResolveErr(id, err.Error())
				}
				hoist.globalsByName[id.Name] = sym
				hoist.globalsById[id.IdId] = sym
				hoist.varOrder = append(hoist.varOrder, id.IdId)
				hoist.varDeclById[id.IdId] = node
				r.setResolved(id, r.refFromSymbol(sym))
			}
		case *parser.FuncDecl:
			sym, err := r.declare(node.Id.Name, SymbolFunc, node.Id.IdId)
			if err != nil {
				return nil, newResolveErr(node.Id, err.Error())
			}
			hoist.globalsByName[node.Id.Name] = sym
			hoist.globalsById[node.Id.IdId] = sym
			hoist.funcOrder = append(hoist.funcOrder, node.Id.IdId)
			hoist.funcDeclById[node.Id.IdId] = node
			r.setResolved(node.Id, r.refFromSymbol(sym))
		}
	}
	return hoist, nil
}
