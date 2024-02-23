package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	BLOCK_SENTINEL = iota

	BLOCK_HTML

	BLOCK_TEMPLATE
	BLOCK_EXPAND
	BLOCK_CONTENT
	BLOCK_USE
	BLOCK_FOR
	BLOCK_END_FOR
	BLOCK_VAR
	BLOCK_OUT_ONLY
)

type Block struct {
	Type int
	Data []string
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

func parseAssignmentExpr(expr string) (string, string, error) {
	return "", "", errors.New("Not implemented")
}

// Parses data between start and end into a special block (non HTML block)
func parseSpecialBlock(data *[]byte, start, end uint64) (*Block, error) {
	// convert to string
	argsString := string((*data)[start:end])
	argsString = strings.ReplaceAll(argsString, "&rsquo;", "'")
	argsString = strings.ReplaceAll(argsString, "&lsquo;", "'")
	re := regexp.MustCompile(`[^\s'\\]+|\\[\\']*|'([^']*?)'`) // keeps quoted string together while splitting on spaces
	// parse out tokens
	tokens := re.FindAllString(argsString, -1)
	if len(tokens) == 0 {
		return nil, errors.New(fmt.Sprintf("Invalid syntax: '%s'", argsString))
	}
	// first token is the name of the block
	strCode := tokens[0]
	blockTypeStr := strings.ToLower(strings.Split(strings.Trim(strCode, " "), " ")[0])
	// parse the rest of it
	switch blockTypeStr {
	case "template":
		var data []string
		// check no template name given
		if len(tokens) < 2 {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `template` tag, no template name given: '%s'", argsString))
		} else {
			data = append(data, tokens[1])
		}
		// check not correct form of expression
		if len(tokens) >= 3 {
			for _, assignExpr := range tokens[2:] {
				key, val, err := parseAssignmentExpr(assignExpr)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Invalid syntax for `template` tag, invalid assignment expression: '%s'", argsString))
				}
				data = append(data, key)
				data = append(data, val)
			}
		}
		return &Block{
			Type: BLOCK_TEMPLATE,
			Data: data,
		}, nil
	case "expand":
		var data []string
		// check no template name given
		if len(tokens) < 2 {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `expand` tag, no template name given: '%s'", argsString))
		} else {
			data = append(data, tokens[1])
		}
		// check not correct form of expression
		if len(tokens) >= 3 {
			for _, assignExpr := range tokens[2:] {
				key, val, err := parseAssignmentExpr(assignExpr)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Invalid syntax for `expand` tag, invalid assignment expression: '%s'", argsString))
				}
				data = append(data, key)
				data = append(data, val)
			}
		}
		return &Block{
			Type: BLOCK_EXPAND,
			Data: data,
		}, nil
	case "content":
		return &Block{
			Type: BLOCK_CONTENT,
		}, nil
	case "use":
		var data []string
		// check invalid length of variables
		if len(tokens) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `use` tag, no variable name given: '%s'", argsString))
		} else {
			data = append(data, tokens[1])
		}
		return &Block{
			Type: BLOCK_USE,
			Data: data,
		}, nil
	case "for":
		var data []string
        // TODO
		blockType = BLOCK_FOR
	case "endfor":
		blockType = BLOCK_END_FOR
	case "var":
		blockType = BLOCK_VAR
	case "outonly":
		blockType = BLOCK_OUT_ONLY
	default:
		return nil, errors.New(fmt.Sprintf("Unrecognized block type '%s'", blockTypeStr))
	}
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
				Data: []string{string((*data)[start:i])},
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
