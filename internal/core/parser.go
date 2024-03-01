package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Checks for '=' sign in the string and returns lhs and rhs
func parseAssignmentExpr(expr string) (string, string, error) {
	index := strings.Index(expr, "=")
	if index == -1 {
		return "", "", errors.New(fmt.Sprintf("Invalid syntax for `var` tag, expected variable assignment: %s", expr))
	}
	lhs := expr[:index]
	rhs := expr[index+1:]
	if len(lhs) == 0 || len(rhs) == 0 {
		return "", "", errors.New(fmt.Sprintf("The LHS and RHS of expression must not be empty: %s", expr))
	}
	return lhs, rhs, nil
}

// Parses data between start and end into a special block (non HTML block)
func parseSpecialBlock(data *[]byte, start, end uint64) (*block, error) {
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
		}
		data = append(data, tokens[1])
		// check not correct form of expression
		if len(tokens) >= 3 {
			for _, assignExpr := range tokens[2:] {
				key, val, err := parseAssignmentExpr(assignExpr)
				if err != nil {
					return nil, err
				}
				data = append(data, key)
				data = append(data, val)
			}
		}
		return &block{
			blockType: blockTypeTemplate,
			data:      data,
		}, nil

	case "expand":
		var data []string
		// check no template name given
		if len(tokens) < 2 {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `expand` tag, no template name given: '%s'", argsString))
		}
		data = append(data, tokens[1])
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
		return &block{
			blockType: blockTypeExapand,
			data:      data,
		}, nil

	case "content":
		return &block{
			blockType: blockTypeContent,
		}, nil

	case "use":
		var data []string
		// check invalid length of variables
		if len(tokens) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `use` tag, no variable name given: '%s'", argsString))
		}
		data = append(data, tokens[1])
		return &block{
			blockType: blockTypeUse,
			data:      data,
		}, nil

	case "for":
		var data []string
		// check invalid length
		if len(tokens) != 4 {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `for` tag: %s", argsString))
		}
		// check for in
		if tokens[2] != "in" {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `for` tag, keyword `in` was expected: %s", argsString))
		}
		data = append(data, tokens[1])
		data = append(data, tokens[3])
		return &block{
			blockType: blockTypeFor,
			data:      data,
		}, nil

	case "endfor":
		return &block{
			blockType: blockTypeEndFor,
		}, nil

	case "var":
		var data []string
		// check invalid length
		if len(tokens) != 2 {
			return nil, errors.New(fmt.Sprintf("Invalid syntax for `var` tag: %s", argsString))
		}
		// parse expression
		key, val, err := parseAssignmentExpr(tokens[1])
		if err != nil {
			return nil, err
		}
		data = append(data, key)
		data = append(data, val)
		return &block{
			blockType: blockTypeVar,
			data:      data,
		}, nil

	case "outonly":
		return &block{
			blockType: blockTypeOutOnly,
		}, nil

	default:
		return nil, errors.New(fmt.Sprintf("Unrecognized block type '%s'", blockTypeStr))
	}
}

// Convert HTML data to chain of blocks
func HTMLToBlocks(data *[]byte) (*blockChain, error) {
	bc := newBlockChain()
	isOpen := false
	var start uint64 = 0
	var dataSize uint64 = uint64(len(*data))
	for i := uint64(0); i < dataSize; i += 1 {
		// opening braces
		if i+1 < dataSize && (*data)[i] == byte('{') && (*data)[i+1] == byte('%') {
			if isOpen {
				return nil, errors.New("Invalid syntax, you cannot nest special blocks")
			}
			if start != i {
				bc.append(&block{
					blockType: blockTypeHtml,
					data:      []string{string((*data)[start:i])},
				})
			}
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
			bc.append(block)
			start = i + 2
			isOpen = false
		}
	}
	// edge case: final html block
	if start != dataSize {
		bc.append(&block{
			blockType: blockTypeHtml,
			data:      []string{string((*data)[start:dataSize])},
		})
	}
	return bc, nil
}
