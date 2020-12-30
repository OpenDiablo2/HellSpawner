package hsnode

import "strings"

const (
	badSep = "\\"
	sep    = "/"
)

type Name = string

type Children map[Name]*Node

type Path = string

type Node struct {
	Parent *Node
	Name
	Children
}

func (n *Node) String() string {
	return n.Name
}

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

func (n *Node) SetParent(p *Node) *Node {
	n.Parent = p

	return n
}

func NewNode(name string) *Node {
	n := &Node{
		Name:     name,
		Children: make(Children),
	}

	return n
}
