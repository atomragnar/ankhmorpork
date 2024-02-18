package filetree

type Visitor interface {
	Visit(node Node)
}

// example usage
//
// type TreeVisitor struct {
// 	Root *Node
// }

// func NewTreeVisitor(node Node) *TreeVisitor {
// 	visitor := &TreeVisitor{
// 		Root: &node,
// 	}

// 	visitor.Visit(node)

// 	return visitor
// }

// func (v *TreeVisitor) Visit(node Node) {
// 	fmt.Println(node.Path)
// 	for _, c := range node.Children {
// 		c.Accept(v)
// 	}
// }
