package deltacoverage

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"

	"golang.org/x/tools/go/ast/astutil"
)

func Instrument(w io.Writer, code string) error {
	fs := token.NewFileSet()
	// we need to parse comments due to CGO code
	file, err := parser.ParseFile(fs, "", code, parser.ParseComments)
	if err != nil {
		return err
	}
	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.FuncDecl:
			// It is possible to have func declaration wihtout body in Go source code
			if x.Body == nil {
				return true
			}
			newCallStmt := &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "fmt",
						},
						Sel: &ast.Ident{
							Name: "Println",
						},
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"instrumentation"`,
						},
					},
				},
			}
			x.Body.List = append([]ast.Stmt{newCallStmt}, x.Body.List...)
			astutil.AddImport(fs, file, "fmt")
		}
		return true
	})

	return printer.Fprint(w, fs, file)
}
