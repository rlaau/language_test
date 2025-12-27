package resolver

import (
	"fmt"

	"github.com/rlaaudgjs5638/langTest/tinygo/parser"
)

type InitStep struct {
	VarId     parser.IdId
	ExprOrNil parser.Expr
	ZeroInit  bool
}

type InitOrder []InitStep

func BuildInitOrder(pkg *parser.PackageAST, table ResolveTable, hoist *HoistInfo) (InitOrder, error) {
	if hoist == nil {
		return nil, fmt.Errorf("hoist info is required")
	}
	// varZeroInit: 초기화되지 않은 모든 VarDecl
	varZeroInit := map[parser.IdId]bool{}
	// varInitFexp: 초기화되었지만 Expr이 Fexp인 VarDecl
	varInitFexp := map[parser.IdId]parser.Expr{}
	// varInitExpr: 초기화되었으며 Expr이 Fexp가 아닌 VarDecl
	varInitExpr := map[parser.IdId]parser.Expr{}

	// VarDecl을 분해해 분류 집합을 만든다.
	for _, decl := range pkg.DeclsOrNil {
		node, ok := decl.(*parser.VarDecl)
		if !ok {
			continue
		}
		if len(node.ExprsOrNil) > 0 && len(node.ExprsOrNil) != len(node.Ids) {
			return nil, fmt.Errorf("var decl lhs/rhs count mismatch")
		}
		if len(node.ExprsOrNil) == 0 {
			for _, id := range node.Ids {
				varZeroInit[id.IdId] = true
			}
			continue
		}
		for i, id := range node.Ids {
			expr := node.ExprsOrNil[i]
			if exprIsFexp(expr) {
				varInitFexp[id.IdId] = expr
				continue
			}
			varInitExpr[id.IdId] = expr
		}
	}

	// varInitExpr가 의존하는 전역 변수/함수/Fexp를 수집한다.
	varToVarDependency := map[parser.IdId]map[parser.IdId]bool{}  //의존하는 모든 초기화 && 값이 Fexp아닌 varDecl
	varToFexpDependency := map[parser.IdId]map[parser.IdId]bool{} // 의존하는 모든 초기화 && 값이 Fexp인 VarDecl
	varToFuncDependency := map[parser.IdId]map[parser.IdId]bool{} // 의존하는 모든 FunclDecl
	for varId, expr := range varInitExpr {
		// depVars는 해당 varId의 expr이 의존하는 모든 varDecl임
		// 이 varDecl은 값이 fexp인 varDecl도 포함. 말 그대로 모든 varDecl임
		// depFunc는 해당 expr이 의존하는 모든 funcDecl임
		depVars, depFuncs, err := collectGlobalRefsExpr(expr, table, hoist)
		if err != nil {
			return nil, err
		}
		for depVar := range depVars {
			//초기화되지 않은 변수는 무조건 최우선으로 init되며, 어떤 의존성도 없으므로
			//굳이 의존성 리스트에 넣을 이유가 없음. 어차피 맨 앞에 위치함
			if varZeroInit[depVar] {
				continue
			}

			if varInitFexp[depVar] != nil {
				addDependency(varToFexpDependency, varId, depVar)
				continue
			}
			if varInitExpr[depVar] != nil {
				addDependency(varToVarDependency, varId, depVar)
			}
		}
		for depFunc := range depFuncs {
			addDependency(varToFuncDependency, varId, depFunc)
		}
	}

	// varInitExpr들에 대해서만 위상정렬한다.
	order, err := topoSortVars(hoist.varIds(), varInitExpr, varToVarDependency)
	if err != nil {
		return nil, err
	}
	// 함수가 어떤 VarDecl에 의존하는지 분석한다.
	funcToVarDep, err := collectFuncGlobalVarDeps(pkg, table, hoist)
	if err != nil {
		return nil, err
	}

	// Fexp VarDecl이 어떤 VarDecl에 의존하는지 분석한다.
	fexpToVarDep, err := collectFexpGlobalVarDeps(varInitFexp, table, hoist)
	if err != nil {
		return nil, err
	}

	// callable(함수/클로저) -> var 의존성 그래프를 합친다.
	callableToVar := map[parser.IdId]map[parser.IdId]bool{}
	for funcId, vars := range funcToVarDep {
		callableToVar[funcId] = vars
	}
	for fexpId, vars := range fexpToVarDep {
		callableToVar[fexpId] = vars
	}

	varToCallable := map[parser.IdId]map[parser.IdId]bool{}
	for varId, funcsDep := range varToFuncDependency {
		for funcId := range funcsDep {
			addDependency(varToCallable, varId, funcId)
		}
	}
	for varId, fexpsDep := range varToFexpDependency {
		for fexpId := range fexpsDep {
			addDependency(varToCallable, varId, fexpId)
		}
	}

	// varInitExpr <-> callable 간 순환 차단
	if err := detectVarCallableCycle(varToCallable, callableToVar, varToVarDependency); err != nil {
		return nil, err
	}

	initOrder := InitOrder{}
	// zero-init은 항상 먼저
	for idId := range varZeroInit {
		initOrder = append(initOrder, InitStep{
			VarId:     idId,
			ExprOrNil: nil,
			ZeroInit:  true,
		})
	}

	// fexp init은 zero-init 다음
	for varId, fexp := range varInitFexp {
		if fexp == nil {
			continue
		}
		initOrder = append(initOrder, InitStep{
			VarId:     varId,
			ExprOrNil: fexp,
		})
	}

	// 나머지 expr init은 위상정렬 순서
	for _, varId := range order {
		initOrder = append(initOrder, InitStep{
			VarId:     varId,
			ExprOrNil: varInitExpr[varId],
		})
	}

	return initOrder, nil
}

func exprIsFexp(expr parser.Expr) bool {
	primary, ok := expr.(*parser.Primary)
	if !ok {
		return false
	}
	if primary.PrimaryKind != parser.ValuePrimary || primary.ValueOrNil == nil {
		return false
	}
	return primary.ValueOrNil.ValueKind == parser.FexpValue
}

// expr이 의존중인 idId를 모두 수합해서 리턴
func collectGlobalRefsExpr(expr parser.Expr, table ResolveTable, hoist *HoistInfo) (map[parser.IdId]bool, map[parser.IdId]bool, error) {
	vars := map[parser.IdId]bool{}
	funcs := map[parser.IdId]bool{}
	if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
		return nil, nil, err
	}
	return vars, funcs, nil
}

// expr이 의존중인 idId를, Expr을 Walk하며 모두 수합하는 함수
func walkExprRefs(expr parser.Expr, table ResolveTable, hoist *HoistInfo, vars, funcs map[parser.IdId]bool) error {
	switch node := expr.(type) {
	case *parser.Unary:
		return walkExprRefs(node.Object, table, hoist, vars, funcs)
	case *parser.Binary:
		if err := walkExprRefs(node.LeftExpr, table, hoist, vars, funcs); err != nil {
			return err
		}
		return walkExprRefs(node.RightExpr, table, hoist, vars, funcs)
	case *parser.Primary:
		return walkPrimaryRefs(node, table, hoist, vars, funcs)
	case *parser.Call:
		if err := walkPrimaryRefs(&node.PrimaryOrNil, table, hoist, vars, funcs); err != nil {
			return err
		}
		for _, args := range node.ArgsList {
			for _, arg := range args {
				if err := walkExprRefs(arg, table, hoist, vars, funcs); err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func walkPrimaryRefs(node *parser.Primary, table ResolveTable, hoist *HoistInfo, vars, funcs map[parser.IdId]bool) error {
	switch node.PrimaryKind {
	case parser.ExprPrimary:
		return walkExprRefs(node.ExprOrNil, table, hoist, vars, funcs)
	case parser.IdPrimary:
		id := *node.IdOrNil
		ref, ok := table[id.IdId]
		if !ok {
			return fmt.Errorf("missing resolve entry for id %s", id.String())
		}
		if ref.Kind != RefGlobal {
			return nil
		}
		sym := hoist.getById(ref.RefIdNodeId)
		if sym == nil {
			return fmt.Errorf("missing hoist entry for id #%d", ref.RefIdNodeId)
		}
		switch sym.kind {
		case SymbolVar:
			vars[sym.idNodeId] = true
		case SymbolFunc:
			funcs[sym.idNodeId] = true
		}
		return nil
	case parser.ValuePrimary:
		if node.ValueOrNil != nil && node.ValueOrNil.ValueKind == parser.FexpValue {
			return walkBlockRefs(node.ValueOrNil.FexpOrNil.Block, table, hoist, vars, funcs)
		}
		return nil
	default:
		return nil
	}
}

func walkBlockRefs(block parser.Block, table ResolveTable, hoist *HoistInfo, vars, funcs map[parser.IdId]bool) error {
	for _, stmt := range block.StmtsOrNil {
		if err := walkStmtRefs(stmt, table, hoist, vars, funcs); err != nil {
			return err
		}
	}
	return nil
}

func walkStmtRefs(stmt parser.Stmt, table ResolveTable, hoist *HoistInfo, vars, funcs map[parser.IdId]bool) error {
	switch node := stmt.(type) {

	case *parser.Assign:
		for _, expr := range node.Exprs {
			if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
				return err
			}
		}
	case *parser.CallStmt:
		return walkExprRefs(&node.Call, table, hoist, vars, funcs)
	case *parser.ShortDecl:
		for _, expr := range node.Exprs {
			if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
				return err
			}
		}
	case *parser.VarDecl:
		for _, expr := range node.ExprsOrNil {
			if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
				return err
			}
		}
	case *parser.FuncDecl:
		return walkBlockRefs(node.Block, table, hoist, vars, funcs)
	case *parser.Return:
		for _, expr := range node.ExprsOrNil {
			if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
				return err
			}
		}
	case *parser.If:
		if node.ShortDeclOrNil != nil {
			for _, expr := range node.ShortDeclOrNil.Exprs {
				if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
					return err
				}
			}
		}
		if err := walkExprRefs(node.Bexp, table, hoist, vars, funcs); err != nil {
			return err
		}
		if err := walkBlockRefs(node.ThenBlock, table, hoist, vars, funcs); err != nil {
			return err
		}
		if node.ElseOrNil != nil {
			if err := walkBlockRefs(*node.ElseOrNil, table, hoist, vars, funcs); err != nil {
				return err
			}
		}
	case *parser.ForBexp:
		if err := walkExprRefs(node.Bexp, table, hoist, vars, funcs); err != nil {
			return err
		}
		return walkBlockRefs(node.Block, table, hoist, vars, funcs)
	case *parser.ForWithAssign:
		for _, expr := range node.ShortDecl.Exprs {
			if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
				return err
			}
		}
		if err := walkExprRefs(node.Bexp, table, hoist, vars, funcs); err != nil {
			return err
		}
		for _, expr := range node.Assign.Exprs {
			if err := walkExprRefs(expr, table, hoist, vars, funcs); err != nil {
				return err
			}
		}
		return walkBlockRefs(node.Block, table, hoist, vars, funcs)
	case *parser.ForRangeAexp:
		if err := walkExprRefs(node.Aexp, table, hoist, vars, funcs); err != nil {
			return err
		}
		return walkBlockRefs(node.Block, table, hoist, vars, funcs)
	case *parser.Block:
		return walkBlockRefs(*node, table, hoist, vars, funcs)
	}
	return nil
}

func collectFuncGlobalVarDeps(pkg *parser.PackageAST, table ResolveTable, hoist *HoistInfo) (map[parser.IdId]map[parser.IdId]bool, error) {
	funcToVar := map[parser.IdId]map[parser.IdId]bool{}
	for _, decl := range pkg.DeclsOrNil {
		fn, ok := decl.(*parser.FuncDecl)
		if !ok {
			continue
		}
		vars := map[parser.IdId]bool{}
		funcs := map[parser.IdId]bool{}
		if err := walkBlockRefs(fn.Block, table, hoist, vars, funcs); err != nil {
			return nil, err
		}
		for varId := range vars {
			addDependency(funcToVar, fn.Id.IdId, varId)
		}
	}
	return funcToVar, nil
}

func collectFexpGlobalVarDeps(varInitFexp map[parser.IdId]parser.Expr, table ResolveTable, hoist *HoistInfo) (map[parser.IdId]map[parser.IdId]bool, error) {
	fexpToVar := map[parser.IdId]map[parser.IdId]bool{}
	for varId, expr := range varInitFexp {
		block := fexpBlockFromExpr(expr)
		if block == nil {
			continue
		}
		vars := map[parser.IdId]bool{}
		funcs := map[parser.IdId]bool{}
		if err := walkBlockRefs(*block, table, hoist, vars, funcs); err != nil {
			return nil, err
		}
		for depVar := range vars {
			addDependency(fexpToVar, varId, depVar)
		}
	}
	return fexpToVar, nil
}

func fexpBlockFromExpr(expr parser.Expr) *parser.Block {
	primary, ok := expr.(*parser.Primary)
	if !ok {
		return nil
	}
	if primary.PrimaryKind != parser.ValuePrimary || primary.ValueOrNil == nil {
		return nil
	}
	if primary.ValueOrNil.ValueKind != parser.FexpValue || primary.ValueOrNil.FexpOrNil == nil {
		return nil
	}
	block := primary.ValueOrNil.FexpOrNil.Block
	return &block
}

func detectVarCallableCycle(varToCallable, callableToVar, varDeps map[parser.IdId]map[parser.IdId]bool) error {
	varClosure := map[parser.IdId]map[parser.IdId]bool{}
	for varId := range varDeps {
		if _, ok := varClosure[varId]; ok {
			continue
		}
		closure, err := buildVarClosure(varId, varDeps, map[parser.IdId]int{})
		if err != nil {
			return err
		}
		varClosure[varId] = closure
	}

	callableReach := map[parser.IdId]map[parser.IdId]bool{}
	for callId, vars := range callableToVar {
		reach := map[parser.IdId]bool{}
		for varId := range vars {
			reach[varId] = true
			if closure, ok := varClosure[varId]; ok {
				for dep := range closure {
					reach[dep] = true
				}
			}
		}
		callableReach[callId] = reach
	}

	for varId, calls := range varToCallable {
		for callId := range calls {
			if callableReach[callId][varId] {
				return fmt.Errorf("cycle detected between var #%d and callable #%d", varId, callId)
			}
		}
	}
	return nil
}

func buildVarClosure(start parser.IdId, varDeps map[parser.IdId]map[parser.IdId]bool, state map[parser.IdId]int) (map[parser.IdId]bool, error) {
	if state[start] == 1 {
		return nil, fmt.Errorf("cycle detected among vars at #%d", start)
	}
	if state[start] == 2 {
		return map[parser.IdId]bool{}, nil
	}
	state[start] = 1
	closure := map[parser.IdId]bool{}
	for dep := range varDeps[start] {
		closure[dep] = true
		sub, err := buildVarClosure(dep, varDeps, state)
		if err != nil {
			return nil, err
		}
		for k := range sub {
			closure[k] = true
		}
	}
	state[start] = 2
	return closure, nil
}

// topoSortVars는
//   - 변수 초기화 간 의존성 그래프를 받아
//   - 위상 정렬된 변수 ID 목록을 반환한다
//   - 순환 의존이 있으면 에러를 반환한다
func topoSortVars(varOrder []parser.IdId, varInitExpr map[parser.IdId]parser.Expr, varToVarDeps map[parser.IdId]map[parser.IdId]bool) ([]parser.IdId, error) {
	// state[id]:
	// 0 = 아직 방문 안 함
	// 1 = 방문 중 (DFS 스택에 있음)(=체인에 올려져 있음)
	// 2 = 방문 완료 (위상정렬 결과에 이미 반영됨)
	state := map[parser.IdId]int{}
	result := []parser.IdId{}

	var visit func(parser.IdId) error
	visit = func(id parser.IdId) error {
		if state[id] == 1 {
			return fmt.Errorf("cycle detected among vars at #%d", id)
		}
		if state[id] == 2 {
			return nil
		}
		// 이 변수를 우선순위 체인에 올림
		state[id] = 1
		// 체인에 올린 상태로 의존성 목록에 대해 의존성 사이클 체크
		// 가장 끄트머리에 있는 의존이, 재귀 함수 하에서 가장 먼저 결과에 append
		for dep := range varToVarDeps[id] {
			if err := visit(dep); err != nil {
				return err
			}
		}
		//사이클 없다면 상태 바꾼 후 result에 포함
		state[id] = 2
		result = append(result, id)
		return nil
	}

	for _, id := range varOrder {
		if varInitExpr[id] == nil {
			continue
		}
		if err := visit(id); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func addDependency[T comparable](deps map[T]map[T]bool, from, to T) {
	if deps[from] == nil {
		deps[from] = map[T]bool{}
	}
	deps[from][to] = true
}
