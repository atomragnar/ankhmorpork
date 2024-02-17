package filetree

import "fmt"

func PrintTree(node *Node, prefix string) {
	entryType := "File"
	if node.Info.IsDir {
		entryType = "Dir"
	}
	fmt.Printf("%s[%s] %s\n", prefix, entryType, node.Path)

	if node.Info.IsDir {
		for _, child := range node.Children {
			PrintTree(child, prefix+"  ")
		}
	}
}
