package filetree

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
)

//type Tree interface {
//	Root() TreeNode
//}

// type NodeInterface interface {
// 	Accept(visitor Visitor) error
// }

type Node struct {
	Path     string
	Info     SerializableFileInfo
	Children []*Node
}

func EncodeNode(node *Node) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(node)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeNode(data []byte) (*Node, error) {
	var node Node
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&node)
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (node *Node) Accept(visitor Visitor) error {
	return visitor.Visit(*node)
}

func (node *Node) AddChild(child *Node) {
	if child == nil {
		node.Children = append(node.Children, child)
	}
}

func New(rootPath string) (*Node, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Path: rootPath,
		Info: NewSerializableFileInfo(info),
	}

	if info.IsDir() {
		entries, err := os.ReadDir(rootPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			childPath := filepath.Join(rootPath, entry.Name())
			childNode, err := New(childPath)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}
