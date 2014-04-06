package template

import (
	"github.com/surullabs/fault"
	"strings"
	"text/template/parse"
)

var check = fault.NewChecker()

// Existence check for variables in templates

// MapData wraps a map passed to the template.Execute function.
type MapData interface {
	Get(string) interface{}
}

type FieldCheck struct {
	Name       string
	Text       string
	LeftDelim  string
	RightDelim string
	Funcs      []map[string]interface{}
}

func CheckMissingFields(tpl FieldCheck, data MapData) (err error) {
	defer check.Recover(&err)
	if tpl.Funcs == nil {
		tpl.Funcs = make([]map[string]interface{}, 0, 0)
	}

	//check.True(tpl.Text != "", "no template text present when checking for missing fields")
	trees := check.Return(parse.Parse(tpl.Name, tpl.Text, tpl.LeftDelim, tpl.RightDelim, tpl.Funcs...)).(map[string]*parse.Tree)
	for _, tree := range trees {
		walkTree(tree.Root, data)
	}
	return
}

func walkTreeNodes(data MapData, nodes ...parse.Node) {
	for _, l := range nodes {
		walkTree(l, data)
	}
}

func exists(data MapData, keys []string) bool {
	var m interface{}
	m = data
	for _, key := range keys {
		mapData, ok := m.(MapData)
		if !ok {
			return false
		}

		m = mapData.Get(key)
		if m == nil {
			return false
		}
	}
	return true

}

func walkTree(node parse.Node, data MapData) {
	switch n := node.(type) {
	case *parse.ActionNode:
		walkTree(n.Pipe, data)
	case *parse.BoolNode:
	case *parse.DotNode:
	case *parse.ChainNode:
	case *parse.NilNode:
	case *parse.NumberNode:
	case *parse.RangeNode:
	case *parse.StringNode:
	case *parse.TemplateNode:
	case *parse.TextNode:
	case *parse.CommandNode:
		walkTreeNodes(data, n.Args...)
	case *parse.FieldNode:
		check.Truef(exists(data, n.Ident), "missing field %s", strings.Join(n.Ident, "."))
	case *parse.IfNode:
		walkTreeNodes(data, n.Pipe, n.List, n.ElseList)
	case *parse.IdentifierNode:
	case *parse.ListNode:
		if n != nil {
			walkTreeNodes(data, n.Nodes...)
		}
	case *parse.PipeNode:
		for _, sub := range n.Cmds {
			walkTree(sub, data)
		}
		for _, sub := range n.Decl {
			walkTree(sub, data)
		}
	case *parse.VariableNode:
	case *parse.WithNode:
		walkTreeNodes(data, n.Pipe, n.List, n.ElseList)
	}
}
