package core

import (
	"errors"
	"fmt"
	"strings"
)

const (
	BLOCK_HTML = iota
	BLOCK_CONTENT
	BLOCK_USE
	BLOCK_DEF
	BLOCK_FOR
	BLOCK_END_FOR
	BLOCK_OUT_ONLY
	BLOCK_EXPAND
	BLOCK_TEMPLATE
	BLOCK_SENTINEL
)

type Block struct {
	Type     int
	Next     *Block
	Prev     *Block
	StartIdx uint64
	EndIdx   uint64
}

type BlockChain struct {
	sentinel *Block
}

func NewBlockChain() *BlockChain {
	// create and return a new block chain
	sentinel := Block{}
	sentinel.Type = BLOCK_SENTINEL
	sentinel.Next = &sentinel
	sentinel.Prev = &sentinel
	return &BlockChain{&sentinel}
}

func ParseBlockType(data *[]byte, start, end uint64) (int, error) {
	// parse the type of block
	// TODO: don't need to make a copy
	code := make([]byte, end-start+1)
	copy(code, (*data)[start:end])
	strCode := string(code)
	blockTypeStr := strings.ToLower(strings.Split(strings.Trim(strCode, " "), " ")[0])
	if blockTypeStr == "template" {
		return BLOCK_TEMPLATE, nil
	} else if blockTypeStr == "expand" {
		return BLOCK_EXPAND, nil
	} else if blockTypeStr == "content" {
		return BLOCK_CONTENT, nil
	} else if blockTypeStr == "use" {
		return BLOCK_USE, nil
	} else if blockTypeStr == "def" {
		return BLOCK_DEF, nil
	} else if blockTypeStr == "for" {
		return BLOCK_FOR, nil
	} else if blockTypeStr == "endfor" {
		return BLOCK_END_FOR, nil
	} else if blockTypeStr == "outonly" {
		return BLOCK_OUT_ONLY, nil
	}
	return -1, errors.New(fmt.Sprintf("Unrecognized block type '%s'", blockTypeStr))
}

func (bc *BlockChain) Append(b *Block) {
	// insert new block at the end of the list
	last := bc.sentinel.Prev
	b.Next = bc.sentinel
	b.Prev = last
	last.Next = b
	bc.sentinel.Prev = b
}

func (bc *BlockChain) AppendLeft(b *Block) {
	// insert new block at the head of the list
	next := bc.sentinel.Next
	b.Next = next
	b.Prev = bc.sentinel
	next.Prev = b
	bc.sentinel.Next = b
}

//func (bc *BlockChain) PopLeft() *Block {
	// remove a block from the start of list, return error if not present
//}

//func (bc *BlockChain) Pop() *Block {
	// remove a block from the end of list, return error if not present
//}

//func (bc *BlockChain) Eject() (*Block, *Block) {
	// remove sentinel and return the head and tail of the list, return error if list is empty
//}
