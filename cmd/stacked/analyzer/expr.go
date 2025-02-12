package analyzer

import "go/ast"

func areExprsEqual(a, b ast.Expr) bool {
	switch a := a.(type) {
	case *ast.Ident:
		b, ok := b.(*ast.Ident)
		return ok && a.Name == b.Name
	case *ast.BasicLit:
		b, ok := b.(*ast.BasicLit)
		return ok && a.Kind == b.Kind && a.Value == b.Value
	case *ast.CompositeLit:
		b, ok := b.(*ast.CompositeLit)
		if !ok || len(a.Elts) != len(b.Elts) {
			return false
		}
		if !areExprsEqual(a.Type, b.Type) {
			return false
		}
		for i := range a.Elts {
			if !areExprsEqual(a.Elts[i], b.Elts[i]) {
				return false
			}
		}
		return false
	case *ast.UnaryExpr:
		b, ok := b.(*ast.UnaryExpr)
		return ok && a.Op == b.Op && areExprsEqual(a.X, b.X)
	case *ast.BinaryExpr:
		b, ok := b.(*ast.BinaryExpr)
		return ok && a.Op == b.Op && areExprsEqual(a.X, b.X) && areExprsEqual(a.Y, b.Y)
	case *ast.CallExpr:
		b, ok := b.(*ast.CallExpr)
		if !ok || len(a.Args) != len(b.Args) {
			return false
		}
		if !areExprsEqual(a.Fun, b.Fun) {
			return false
		}
		for i := range a.Args {
			if !areExprsEqual(a.Args[i], b.Args[i]) {
				return false
			}
		}
		return true
	case *ast.ParenExpr:
		b, ok := b.(*ast.ParenExpr)
		return ok && areExprsEqual(a.X, b.X)
	case *ast.SelectorExpr:
		b, ok := b.(*ast.SelectorExpr)
		return ok && a.Sel.Name == b.Sel.Name && areExprsEqual(a.X, b.X)
	case *ast.IndexExpr:
		b, ok := b.(*ast.IndexExpr)
		return ok && areExprsEqual(a.X, b.X) && areExprsEqual(a.Index, b.Index)
	case *ast.IndexListExpr:
		b, ok := b.(*ast.IndexListExpr)
		if !ok || len(a.Indices) != len(b.Indices) {
			return false
		}
		if !areExprsEqual(a.X, b.X) {
			return false
		}
		for i := range a.Indices {
			if !areExprsEqual(a.Indices[i], b.Indices[i]) {
				return false
			}
		}
		return true
	case *ast.SliceExpr:
		b, ok := b.(*ast.SliceExpr)
		return ok && a.Slice3 == b.Slice3 &&
			areExprsEqual(a.X, b.X) &&
			areExprsEqual(a.Low, b.Low) &&
			areExprsEqual(a.High, b.High) &&
			areExprsEqual(a.Max, b.Max)
	case *ast.KeyValueExpr:
		b, ok := b.(*ast.KeyValueExpr)
		return ok && areExprsEqual(a.Key, b.Key) && areExprsEqual(a.Value, b.Value)
	case *ast.TypeAssertExpr:
		b, ok := b.(*ast.TypeAssertExpr)
		return ok && areExprsEqual(a.X, b.X) && areExprsEqual(a.Type, b.Type)
	case *ast.StarExpr:
		b, ok := b.(*ast.StarExpr)
		return ok && areExprsEqual(a.X, b.X)
	}

	return false
}
