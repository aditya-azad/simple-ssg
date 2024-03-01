package core

import (
	"bytes"
	"strings"
)

const (
	blockTypeSentinel = iota
	blockTypeHtml
	blockTypeTemplate
	blockTypeExapand
	blockTypeContent
	blockTypeUse
	blockTypeFor
	blockTypeEndFor
	blockTypeVar
	blockTypeOutOnly
)

type block struct {
	blockType int
	data      []string
	next      *block
	prev      *block
}

type blockChain struct {
	sentinel *block
}

// Create and return a new block chain
func newBlockChain() *blockChain {
	sentinel := block{}
	sentinel.blockType = blockTypeSentinel
	sentinel.next = &sentinel
	sentinel.prev = &sentinel
	return &blockChain{&sentinel}
}

// Insert new block at the end of the list
func (bc *blockChain) append(b *block) {
	last := bc.sentinel.prev
	b.next = bc.sentinel
	b.prev = last
	last.next = b
	bc.sentinel.prev = b
}

// Insert new block at the head of the list
func (bc *blockChain) appendLeft(b *block) {
	next := bc.sentinel.next
	b.next = next
	b.prev = bc.sentinel
	next.prev = b
	bc.sentinel.next = b
}

// Convert the chain to a string representation
func (bc *blockChain) toString(displayData bool) string {
	curr := bc.sentinel
	typeMap := map[int]string{
		blockTypeSentinel: "SENTINEL",
		blockTypeHtml:     "HTML",
		blockTypeTemplate: "TEMPLATE",
		blockTypeExapand:  "EXPAND",
		blockTypeContent:  "CONTENT",
		blockTypeUse:      "USE",
		blockTypeFor:      "FOR",
		blockTypeEndFor:   "END FOR",
		blockTypeVar:      "VAR",
		blockTypeOutOnly:  "OUT ONLY",
	}
	var buffer bytes.Buffer
	for {
		// print type
		buffer.WriteString(typeMap[curr.blockType])
		// print data
		if displayData && curr.data != nil {
			buffer.WriteString(" (")
			buffer.WriteString(strings.ReplaceAll(strings.Join(curr.data, ", "), "\r\n", ""))
			buffer.WriteString(")")
		}
		buffer.WriteString("\n")
		curr = curr.next
		if curr.blockType == blockTypeSentinel {
			break
		}
	}
	return buffer.String()
}
