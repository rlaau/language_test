package resolver

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

type HoistInfo struct {
	globalsByName map[string]*Symbol
	globalsById   map[parser.IdId]*Symbol
	// 선언 순서대로 수집
	// var간 초기화 순서는 추후 initOrder에서 보장함
	varOrder  []parser.IdId
	funcOrder []parser.IdId
}

func newHoistInfo() *HoistInfo {
	return &HoistInfo{
		globalsByName: map[string]*Symbol{},
		globalsById:   map[parser.IdId]*Symbol{},
	}
}

func (h *HoistInfo) get(name string) *Symbol {
	return h.globalsByName[name]
}

func (h *HoistInfo) getById(id parser.IdId) *Symbol {
	return h.globalsById[id]
}

func (h *HoistInfo) varIds() []parser.IdId {
	return h.varOrder
}

func (h *HoistInfo) funcIds() []parser.IdId {
	return h.funcOrder
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
			r.setResolved(node.Id, r.refFromSymbol(sym))
		}
	}
	return hoist, nil
}
