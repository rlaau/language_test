package resolver

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

func Resolve(pkg *parser.PackageAST) (ResolveTable, *HoistInfo, InitOrder, map[string]int, error) {
	rs := NewResolver()
	table, hoist, err := rs.ResolvePackage(pkg)
	if err != nil {
		return table, hoist, nil, nil, err
	}
	order, ierr := BuildInitOrder(table, hoist)
	if ierr != nil {
		return table, hoist, order, nil, ierr
	}
	return table, hoist, order, copyBuiltins(rs.builtins), nil
}

func copyBuiltins(src map[string]int) map[string]int {
	dst := make(map[string]int, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
