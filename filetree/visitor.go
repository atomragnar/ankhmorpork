package filetree

type Visitor interface {
	Visit(node Node) error
}

type TreeVisitor struct {
	Root *Node
}

func (v *TreeVisitor) Visit(node Node) error {
	// ...
	err := node.Accept(v)
	if err != nil {
		return err
	}
	return nil
}
