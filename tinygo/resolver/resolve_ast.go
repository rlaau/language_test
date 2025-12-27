package resolver

import (
	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

type InitOrder []parser.IdId

func (r *Resolver) ResolvePackage(pkg *parser.PackageAST) (ResolveTable, InitOrder, error) {
	// 패키지 레벨의 선언은 호이스팅함
	hoist, err := r.collectPackageDecls(pkg)
	if err != nil {
		return r.table, nil, err
	}
	for _, decl := range pkg.DeclsOrNil {
		//호이스팅된 정보를 가지고 DFS식 리졸빙 시작
		if err := r.resolveDeclWithHoist(decl, hoist); err != nil {
			return r.table, nil, err
		}
	}
	return r.table, nil, nil
}

// 패키지 레벨의 선언은 호이스팅됨
func (r *Resolver) resolveDeclWithHoist(decl parser.Decl, hoist *HoistInfo) error {
	switch node := decl.(type) {
	case *parser.VarDecl:
		return r.resolveHoistedVarDecl(node, hoist)
	case *parser.FuncDecl:
		return r.resolveHoistedFuncDecl(node, hoist)
	default:
		return nil
	}
}

func (r *Resolver) resolveStmt(stmt parser.Stmt) error {
	switch node := stmt.(type) {
	case *parser.Assign:
		return r.resolveAssign(node)
	case *parser.CallStmt:
		return r.resolveCall(node.Call)
	case *parser.ShortDecl:
		return r.resolveShortDecl(node)
	case *parser.VarDecl:
		return r.resolveVarDecl(node)
	case *parser.FuncDecl:
		return r.resolveFuncDecl(node)
	case *parser.Return:
		return r.resolveReturn(node)
	case *parser.If:
		return r.resolveIf(node)
	case *parser.ForBexp:
		return r.resolveForBexp(node)
	case *parser.ForWithAssign:
		return r.resolveForWithAssign(node)
	case *parser.ForRangeAexp:
		return r.resolveForRangeAexp(node)
	case *parser.Block:
		// 그냥 블록 시엔 새 스코프
		return r.resolveBlock(*node, false)
	default:
		return nil
	}
}

func (r *Resolver) resolveExpr(expr parser.Expr) error {
	switch node := expr.(type) {
	case *parser.Unary:
		return r.resolveExpr(node.Object)
	case *parser.Binary:
		if err := r.resolveExpr(node.LeftExpr); err != nil {
			return err
		}
		return r.resolveExpr(node.RightExpr)
	case *parser.Primary:
		return r.resolvePrimary(node)
	case *parser.Call:
		return r.resolveCall(*node)
	default:
		return nil
	}
}

func (r *Resolver) resolvePrimary(node *parser.Primary) error {
	switch node.PrimaryKind {
	case parser.ExprPrimary:
		return r.resolveExpr(node.ExprOrNil)
	case parser.IdPrimary:
		ref, err := r.resolveID(*node.IdOrNil)
		if err != nil {
			return err
		}
		r.setResolved(*node.IdOrNil, ref)
		return nil
	case parser.ValuePrimary:
		return r.resolveValueForm(node.ValueOrNil)
	default:
		return nil
	}
}

func (r *Resolver) resolveValueForm(v *parser.ValueForm) error {
	if v == nil {
		return nil
	}
	if v.ValueKind == parser.FexpValue {
		return r.resolveFexp(v.FexpOrNil)
	}
	return nil
}

func (r *Resolver) resolveFexp(f *parser.Fexp) error {
	r.pushScope()
	defer r.popScope()
	// fexp: params이전부터 새 스코프
	for _, param := range f.ParamsOrNil {
		sym, err := r.declare(param.Id.Name, SymbolParam, param.Id.IdId)
		if err != nil {
			return newResolveErr(param.Id, err.Error())
		}
		r.setResolved(param.Id, r.refFromSymbol(sym))
	}
	// true에 의해 block은 스코프 재사용
	return r.resolveBlock(f.Block, true)
}

func (r *Resolver) resolveCall(call parser.Call) error {
	if err := r.resolvePrimary(&call.PrimaryOrNil); err != nil {
		return err
	}
	for _, args := range call.ArgsList {
		for _, expr := range args {
			if err := r.resolveExpr(expr); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Resolver) resolveVarDecl(node *parser.VarDecl) error {
	// LHS, RHS 검사
	if len(node.ExprsOrNil) > 0 && len(node.Ids) != len(node.ExprsOrNil) {
		return newResolveErr(firstId(node.Ids), "var decl lhs/rhs count mismatch")
	}
	// 우변 먼저 resolve (재귀적 정의 차단)
	for _, expr := range node.ExprsOrNil {
		if err := r.resolveExpr(expr); err != nil {
			return err
		}
	}
	//좌변 resolve && 좌변은 모두 새 변수
	for _, id := range node.Ids {
		sym, err := r.declare(id.Name, SymbolVar, id.IdId)
		if err != nil {
			return newResolveErr(id, err.Error())
		}
		r.setResolved(id, r.refFromSymbol(sym))
	}
	return nil
}

func (r *Resolver) resolveHoistedVarDecl(node *parser.VarDecl, hoist *HoistInfo) error {
	if len(node.ExprsOrNil) > 0 && len(node.Ids) != len(node.ExprsOrNil) {
		return newResolveErr(firstId(node.Ids), "var decl lhs/rhs count mismatch")
	}
	//우변 리졸브
	for _, expr := range node.ExprsOrNil {
		if err := r.resolveExpr(expr); err != nil {
			return err
		}
	}
	//패키지 레벨의 호이스팅된 VarDecl은
	// 잘 호이스팅 되었나 검사만 함
	for _, id := range node.Ids {
		sym := hoist.get(id.Name)
		if sym == nil {
			return newResolveErr(id, "hoisted symbol not found")
		}
	}
	return nil
}

func (r *Resolver) resolveFuncDecl(node *parser.FuncDecl) error {
	// 좌변 먼저 리졸브 (재귀 허용)
	sym, err := r.declare(node.Id.Name, SymbolFunc, node.Id.IdId)
	if err != nil {
		return newResolveErr(node.Id, err.Error())
	}
	r.setResolved(node.Id, r.refFromSymbol(sym))

	//이후 우변 리졸브
	r.pushScope()
	defer r.popScope()
	for _, param := range node.ParamsOrNil {
		psym, perr := r.declare(param.Id.Name, SymbolParam, param.Id.IdId)
		if perr != nil {
			return newResolveErr(param.Id, perr.Error())
		}
		r.setResolved(param.Id, r.refFromSymbol(psym))
	}
	return r.resolveBlock(node.Block, true)
}

func (r *Resolver) resolveHoistedFuncDecl(node *parser.FuncDecl, hoist *HoistInfo) error {

	// 좌변 먼저 리졸브
	// But 패키지 레벨 호이스팅된 함수는
	// 이미 수집되었으므로, 잘 수집되었나 검사만 하.ㅁ
	sym := hoist.get(node.Id.Name)
	if sym == nil {
		return newResolveErr(node.Id, "hoisted symbol not found")
	}

	r.pushScope()
	defer r.popScope()

	//우변 리졸브
	for _, param := range node.ParamsOrNil {
		psym, perr := r.declare(param.Id.Name, SymbolParam, param.Id.IdId)
		if perr != nil {
			return newResolveErr(param.Id, perr.Error())
		}
		r.setResolved(param.Id, r.refFromSymbol(psym))
	}
	return r.resolveBlock(node.Block, true)
}

func (r *Resolver) resolveAssign(node *parser.Assign) error {
	//LHS, RHS검사
	if len(node.Ids) != len(node.Exprs) {
		return newResolveErr(firstId(node.Ids), "assign lhs/rhs count mismatch")
	}
	// 우변 먼저 리졸브
	for _, expr := range node.Exprs {
		if err := r.resolveExpr(expr); err != nil {
			return err
		}
	}
	// 좌변 리졸브
	for _, id := range node.Ids {
		ref, err := r.resolveID(id)
		if err != nil {
			return err
		}
		//빌트인엔 할당 불가
		if ref.Kind == RefBuiltin {
			return newResolveErr(id, "cannot assign to builtin")
		}
		r.setResolved(id, ref)
	}
	return nil
}

func (r *Resolver) resolveShortDecl(node *parser.ShortDecl) error {
	//ShortDecl도 LHS, RHS검사
	if len(node.Ids) != len(node.Exprs) {
		return newResolveErr(firstId(node.Ids), "short decl lhs/rhs count mismatch")
	}
	// 우변 먼저 리졸브
	for _, expr := range node.Exprs {
		if err := r.resolveExpr(expr); err != nil {
			return err
		}
	}
	newCount := 0
	for _, id := range node.Ids {
		// 스코프에 존재 시 할당으로 처리
		if sym, ok := r.currentScope.symbols[id.Name]; ok {
			ref, err := r.resolveID(id)
			if err != nil {
				return err
			}
			r.setResolved(id, ref)
			// 빌트인은 할당 불가
			if sym.kind == SymbolBuiltin {
				return newResolveErr(id, "cannot assign to builtin")
			}
			continue
		}
		// 스코프에 없다면 선언으로 처리
		sym, err := r.declare(id.Name, SymbolVar, id.IdId)
		if err != nil {
			return newResolveErr(id, err.Error())
		}
		newCount++
		r.setResolved(id, r.refFromSymbol(sym))
	}
	// ShortDecl은 적어도 하나의 새 변수 필요
	if newCount == 0 {
		return newResolveErr(firstId(node.Ids), "short decl requires at least one new variable")
	}
	return nil
}

func (r *Resolver) resolveReturn(node *parser.Return) error {
	for _, expr := range node.ExprsOrNil {
		if err := r.resolveExpr(expr); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveIf(node *parser.If) error {
	r.pushScope()
	defer r.popScope()

	if node.ShortDeclOrNil != nil {
		if err := r.resolveShortDecl(node.ShortDeclOrNil); err != nil {
			return err
		}
	}
	if err := r.resolveExpr(node.Bexp); err != nil {
		return err
	}

	//if는 block에서 새 스코프 생성
	// if의 shortDecl에서 변수 선언 후, if의 block에서 해당 변수 셰도잉 가능
	if err := r.resolveBlock(node.ThenBlock, false); err != nil {
		return err
	}
	if node.ElseOrNil != nil {
		if err := r.resolveBlock(*node.ElseOrNil, false); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveForBexp(node *parser.ForBexp) error {
	r.pushScope()
	defer r.popScope()

	if err := r.resolveExpr(node.Bexp); err != nil {
		return err
	}
	// for역시 블록과 스코프 분리
	return r.resolveBlock(node.Block, false)
}

func (r *Resolver) resolveForWithAssign(node *parser.ForWithAssign) error {
	r.pushScope()
	defer r.popScope()

	if err := r.resolveShortDecl(&node.ShortDecl); err != nil {
		return err
	}
	if err := r.resolveExpr(node.Bexp); err != nil {
		return err
	}
	if err := r.resolveAssign(&node.Assign); err != nil {
		return err
	}
	// for역시 블록과 스코프 분리
	return r.resolveBlock(node.Block, false)
}

func (r *Resolver) resolveForRangeAexp(node *parser.ForRangeAexp) error {
	r.pushScope()
	defer r.popScope()

	if err := r.resolveExpr(node.Aexp); err != nil {
		return err
	}
	// for역시 블록과 스코프 분리
	return r.resolveBlock(node.Block, false)
}

func (r *Resolver) resolveBlock(block parser.Block, reuseCurrent bool) error {
	// 스코프 합쳐달란 요청이 있었다면
	// 기존의 스코프를 재사용해서 스코프 합침.
	// 아니라면, 새 스코프 push, pop
	if !reuseCurrent {
		r.pushScope()
		defer r.popScope()
	}
	for _, stmt := range block.StmtsOrNil {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) refFromSymbol(sym *Symbol) ResolvedRef {
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
		ref.Slot = r.builtins[sym.name]
	case SymbolFunc, SymbolVar:
		if sym.scope == r.global {
			ref.Kind = RefGlobal
			ref.Distance = 0
		}
	}
	return ref
}

func (r *Resolver) setResolved(id parser.Id, ref ResolvedRef) {
	r.table[id.IdId] = ref
}

func firstId(ids []parser.Id) parser.Id {
	if len(ids) > 0 {
		return ids[0]
	}
	return parser.Id{}
}
