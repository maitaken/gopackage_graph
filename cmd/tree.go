package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var targetModule string

type PackageNode struct {
	childNodes []*PackageNode
	pkgName    string
	indent     int
}

func (n *PackageNode) String() string {
	var b strings.Builder
	recursiveString(n, &b)
	return b.String()
}

func recursiveString(n *PackageNode, b *strings.Builder) {
	fmt.Fprintf(b, "|%s: %s\n", strings.Repeat("--", n.indent), n.pkgName)
	for _, nextNode := range n.childNodes {
		recursiveString(nextNode, b)
	}
}

func (n *PackageNode) Sort() {
	if n.childNodes != nil {
		sort.Slice(n.childNodes, func(i, j int) bool {
			return n.childNodes[i].pkgName < n.childNodes[j].pkgName
		})
	}

	for _, nextNode := range n.childNodes {
		nextNode.Sort()
	}
}

func buildTree(parentNode *PackageNode, pkgs []*packages.Package) {
	childNodes := make([]*PackageNode, 0)

	for _, pkg := range pkgs {
		childNode := &PackageNode{
			pkgName: trippedPkgName(pkg.PkgPath),
			indent:  parentNode.indent + 1,
		}
		childNodes = append(childNodes, childNode)

		nextImports := make([]*packages.Package, 0)
		for _, imp := range pkg.Imports {
			if imp.Module != nil && imp.Module.Path == targetModule {
				nextImports = append(nextImports, imp)
			}
		}
		buildTree(childNode, nextImports)
	}

	parentNode.childNodes = childNodes
}

func buildPackageTree(modPath, packPath string) (*PackageNode, error) {
	baseDir, _ := filepath.Split(packPath)
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax | packages.NeedModule,
		Dir:  modPath,
	}

	pkgs, err := packages.Load(cfg, baseDir)
	if err != nil {
		return nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("failed to load packages")
	}

	targetModule = pkgs[0].Module.Path

	root := &PackageNode{
		indent: 0,
	}

	buildTree(root, pkgs)
	root.Sort()

	return root, nil
}

func trippedPkgName(pkgName string) string {
	return strings.TrimPrefix(pkgName, targetModule+"/")
}
