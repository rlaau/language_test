package resolver

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

type HoistInfo struct {
	globals map[string]*Symbol
}

func newHoistInfo() *HoistInfo {
	return &HoistInfo{
		globals: map[string]*Symbol{},
	}
}

func (h *HoistInfo) get(name string) *Symbol {
	return h.globals[name]
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
				hoist.globals[id.Name] = sym
				r.setResolved(id, r.refFromSymbol(sym))
			}
		case *parser.FuncDecl:
			sym, err := r.declare(node.Id.Name, SymbolFunc, node.Id.IdId)
			if err != nil {
				return nil, newResolveErr(node.Id, err.Error())
			}
			hoist.globals[node.Id.Name] = sym
			r.setResolved(node.Id, r.refFromSymbol(sym))
		}
	}
	return hoist, nil
}
