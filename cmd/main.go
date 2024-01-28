package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <module path> <package path>\n", os.Args[0])
		fmt.Printf("  get args: %v\n", os.Args)
		os.Exit(1)
	}

	node, err := buildPackageTree(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("failed to build package tree: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(node)
}
