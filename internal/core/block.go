package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	BLOCK_HTML = iota
	BLOCK_CONTENT
	BLOCK_USE
	BLOCK_FOR
	BLOCK_END_FOR
	BLOCK_OUT_ONLY
	BLOCK_EXPAND
	BLOCK_TEMPLATE
	BLOCK_SENTINEL
)

type Block struct {
	Type int
	Args []string
	Next *Block
	Prev *Block
}

type BlockChain struct {
	sentinel *Block
}

// Create and return a new block chain
func NewBlockChain() *BlockChain {
	sentinel := Block{}
	sentinel.Type = BLOCK_SENTINEL
	sentinel.Next = &sentinel
	sentinel.Prev = &sentinel
	return &BlockChain{&sentinel}
}

// Insert new block at the end of the list
func (bc *BlockChain) Append(b *Block) {
	last := bc.sentinel.Prev
	b.Next = bc.sentinel
	b.Prev = last
	last.Next = b
	bc.sentinel.Prev = b
}

// Insert new block at the head of the list
func (bc *BlockChain) AppendLeft(b *Block) {
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

// Parses data between start and end into a special block (non HTML block)
func parseSpecialBlock(data *[]byte, start, end uint64) (*Block, error) {
	argsString := string((*data)[start:end])
	argsString = strings.ReplaceAll(argsString, "&rsquo;", "'")
	argsString = strings.ReplaceAll(argsString, "&lsquo;", "'")
	re := regexp.MustCompile(`[^\s'\\]+|\\[\\']*|'([^']*?)'`) // keeps quoted string together while splitting on spaces
	tokens := re.FindAllString(argsString, -1)
	// parse the type of block
	strCode := tokens[0]
	blockTypeStr := strings.ToLower(strings.Split(strings.Trim(strCode, " "), " ")[0])
	blockType := -1
	switch blockTypeStr {
	case "template":
		blockType = BLOCK_TEMPLATE
	case "expand":
		blockType = BLOCK_EXPAND
	case "content":
		blockType = BLOCK_CONTENT
	case "use":
		blockType = BLOCK_USE
	case "for":
		blockType = BLOCK_FOR
	case "endfor":
		blockType = BLOCK_END_FOR
	case "outonly":
		blockType = BLOCK_OUT_ONLY
	default:
		return nil, errors.New(fmt.Sprintf("Unrecognized block type '%s'", blockTypeStr))
	}
	// parse args
	for _, x := range tokens {
		fmt.Printf("`%s`,", x)
	}
	fmt.Println()
	return &Block{
		Type: blockType,
		Args: tokens[1:],
	}, nil
}

// Convert HTML data to chain of blocks
func HTMLToBlocks(data *[]byte) (*BlockChain, error) {
	bc := NewBlockChain()
	isOpen := false
	var start uint64 = 0
	var dataSize uint64 = uint64(len(*data))
	for i := uint64(0); i < dataSize; i += 1 {
		// opening braces
		if i+1 < dataSize && (*data)[i] == byte('{') && (*data)[i+1] == byte('%') {
			if isOpen {
				return nil, errors.New("Invalid syntax, you cannot nest special blocks")
			}
			bc.Append(&Block{
				Type: BLOCK_HTML,
				Args: []string{string((*data)[start:i])},
			})
			start = i + 2
			isOpen = true
		}
		// closing braces
		if i+1 < dataSize && (*data)[i] == byte('%') && (*data)[i+1] == byte('}') {
			if !isOpen {
				return nil, errors.New("Invalid syntax, you cannot close a unopened block")
			}
			block, err := parseSpecialBlock(data, start, i)
			if err != nil {
				return nil, err
			}
			bc.Append(block)
			start = i + 2
			isOpen = false
		}
	}
	return bc, nil
}
