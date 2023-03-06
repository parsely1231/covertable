package main

import (
	"errors"
	"path/filepath"

	"golang.org/x/tools/cover"
)

type profileNode interface {
	FilePath() string
	TotalStatements() int
	TestedStatements() int
	AddChild(node profileNode) error
	Children() []profileNode
}

type profileLeaf struct {
	filePath         string
	totalStatements  int
	testedStatements int
}

func (p *profileLeaf) FilePath() string {
	return p.filePath
}

func (p *profileLeaf) TotalStatements() int {
	return p.totalStatements
}

func (p *profileLeaf) TestedStatements() int {
	return p.testedStatements
}

func (*profileLeaf) Children() []profileNode {
	return nil
}

func (*profileLeaf) AddChild(node profileNode) error {
	return errors.New("leaf can't add child")
}

func newProfileLeafFromCoverProfile(profile cover.Profile) *profileLeaf {
	totalStatements := 0
	testedStatements := 0
	for _, b := range profile.Blocks {
		totalStatements += b.NumStmt
		if b.Count >= 1 {
			testedStatements += b.NumStmt
		}
	}

	return &profileLeaf{
		filePath:         profile.FileName,
		totalStatements:  totalStatements,
		testedStatements: testedStatements,
	}
}

type profileInner struct {
	filePath              string
	children              []profileNode
	totalStatementsCache  *int
	testedStatementsCache *int
}

func (p *profileInner) FilePath() string {
	return p.filePath
}

func (p *profileInner) TotalStatements() int {
	if p.totalStatementsCache != nil {
		return *p.totalStatementsCache
	}

	total := 0
	for _, node := range p.children {
		total += node.TotalStatements()
	}
	p.totalStatementsCache = &total
	return total
}

func (p *profileInner) TestedStatements() int {
	if p.testedStatementsCache != nil {
		return *p.testedStatementsCache
	}

	testedTotal := 0
	for _, node := range p.children {
		testedTotal += node.TestedStatements()
	}
	p.testedStatementsCache = &testedTotal
	return testedTotal
}

func (p *profileInner) Children() []profileNode {
	return p.children
}

func (p *profileInner) AddChild(node profileNode) error {
	p.children = append(p.children, node)
	return nil
}

func newProfileInnerFromChild(child profileNode) *profileInner {
	return &profileInner{
		filePath: filepath.Dir(child.FilePath()),
		children: []profileNode{child},
	}
}
