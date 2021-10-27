package db

import (
	"html/template"
	"text/template/parse"
)

// Вытаскивает все параметры из шаблона
func ListTemplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root)
}

func listNodeFields(node parse.Node) []string {
	var res []string
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = append(res, listNodeFields(n)...)
		}
	}
	return res
}
