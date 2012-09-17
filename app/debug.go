package app

import (
	"exp/html"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func debug(w http.ResponseWriter, r *http.Request) {
	node, err := html.Parse(strings.NewReader(manageChecklistsHtml))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	debugNodeFull(w, node, 0)
}

func debugNodeFull(w io.Writer, node *html.Node, n int) {
	if node == nil {
		return
	}
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, " ")
	}
	debugNode(w, node)
	debugNodeFull(w, node.FirstChild, n+2)
	debugNodeFull(w, node.NextSibling, n)
}

func debugNode(w io.Writer, node *html.Node) {
	fmt.Fprintf(
		w,
		"%s <%s>: %s\n",
		nodeTypeToString(node.Type),
		node.DataAtom.String(),
		node.Data,
	)
}

func nodeTypeToString(nodeType html.NodeType) string {
	switch nodeType {
	case html.ErrorNode:
		return "Error"
	case html.TextNode:
		return "Text"
	case html.DocumentNode:
		return "Document"
	case html.ElementNode:
		return "Element"
	case html.CommentNode:
		return "Comment"
	case html.DoctypeNode:
		return "Doctype"
	}
	return "<UnknownType>"
}
