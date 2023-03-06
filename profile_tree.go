package main

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/cover"
)

type profileTree struct {
	root  *profileInner
	nodes map[string]profileNode
}

func makeProfileTree(rootPath string) profileTree {
	root := &profileInner{filePath: rootPath}
	tree := profileTree{
		root:  root,
		nodes: make(map[string]profileNode),
	}
	tree.nodes[rootPath] = root
	return tree
}

func (t profileTree) AddProfile(profile cover.Profile) error {
	profileLeaf := newProfileLeafFromCoverProfile(profile)
	if !strings.HasPrefix(profileLeaf.FilePath(), t.root.FilePath()) {
		return errors.Errorf("profile does not contain root path %v", profileLeaf)
	}

	err := t.connect(profileLeaf)
	if err != nil {
		return err
	}

	return nil
}

func (t profileTree) connect(child profileNode) error {
	parentPath := filepath.Dir(child.FilePath())
	if parent, exist := t.nodes[parentPath]; exist {
		err := parent.AddChild(child)
		if err != nil {
			return err
		}

		t.nodes[child.FilePath()] = child
		return nil
	}

	parent := newProfileInnerFromChild(child)
	err := t.connect(parent)
	if err != nil {
		return err
	}

	t.nodes[child.FilePath()] = child
	return nil
}

func (t profileTree) Walk() []profileNode {
	nodes := make([]profileNode, 0)
	nodes = append(nodes, t.root)

	todo := make([]profileNode, 0)
	todo = append(todo, t.root)
	for len(todo) > 0 {
		node := todo[0]
		todo = todo[1:]

		nodes = append(nodes, node.Children()...)
		todo = append(todo, node.Children()...)
	}

	return nodes
}
