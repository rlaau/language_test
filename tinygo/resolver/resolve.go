package resolver

import "github.com/rlaaudgjs5638/langTest/tinygo/parser"

func Resolve(pkg *parser.PackageAST) (ResolveTable, *HoistInfo, InitOrder, error) {
	rs := NewResolver()
	table, hoist, err := rs.ResolvePackage(pkg)
	if err != nil {
		return table, hoist, nil, err
	}
	order, ierr := BuildInitOrder(table, hoist)
	if ierr != nil {
		return table, hoist, order, ierr
	}
	return table, hoist, order, nil
}
