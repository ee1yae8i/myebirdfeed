package app

import (
	"exp/html"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func debug(w http.ResponseWriter, r *http.Request) {
	rootNode, err := html.Parse(strings.NewReader(manageChecklistsHtml))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	checkListIds := make(map[string]bool)
	for node, depth := rootNode, 0; node != nil; node, depth = nextNode(node, depth) {
		for i := 0; i < depth; i++ {
			fmt.Fprintf(w, "  ")
		}
		debugNode(w, node)
		checkListId := findCheckListId(node)
		if checkListId != nil {
			checkListIds[checkListId] = true
		}
	}
}

// Returns the next node in a depth-first traversal.
func nextNode(node *html.Node, depth int) (*html.Node, int) {
	if node.FirstChild != nil {
		return node.FirstChild, depth+1
	}
	if node.NextSibling != nil {
		return node.NextSibling, depth
	}
	for node = node.Parent; node != nil; node = node.Parent {
		depth -= 1
		if node.NextSibling != nil {
			return node.NextSibling, depth
		}
	}
	return nil, 0
}

func debugNode(w io.Writer, node *html.Node) {
	fmt.Fprintf(w, "%s ", nodeTypeToString(node.Type))
	if node.Type == html.ElementNode {
		fmt.Fprintf(w, "<%s", node.DataAtom.String())
		for _, attr := range node.Attr {
			fmt.Fprintf(w, " %s=\"%s\"", attr.Key, attr.Val)
		}
		fmt.Fprintf(w, ">")
	}

	fmt.Fprintf(w, ": %s\n", node.Data)
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
