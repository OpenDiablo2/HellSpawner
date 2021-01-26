// Pakcage hsnode contains game nodes
package hsnode

import "strings"

const (
	badSep = "\\"
	sep    = "/"
)

// Name represents node's name
type Name = string

// Children represents represents subnodes
type Children map[Name]*Node

// Path represents node path
type Path = string

// Node represents node
type Node struct {
	Parent *Node
	Name
	Children
}

// NewNode creates a new node
func NewNode(name string) *Node {
	n := &Node{
		Name:     name,
		Children: make(Children),
	}

	return n
}

// String represents node's name
func (n *Node) String() string {
	return n.Name
}

// Insert inserts a new path to node
func (n *Node) Insert(p Path) *Node {
	p = strings.ToLower(p)
	p = strings.ReplaceAll(p, badSep, sep)

	parts := strings.Split(p, sep)

	if len(parts) > 1 {
		part := parts[0]
		remainder := strings.Join(parts[1:], sep)

		child, found := n.Children[part]
		if !found {
			n.Children[part] = NewNode(part)
			child = n.Children[part]
		}

		child.Insert(remainder).SetParent(n)
	} else {
		n.Name = p
	}

	return n
}

// SetParent sets node's parent
func (n *Node) SetParent(p *Node) *Node {
	n.Parent = p

	return n
}
