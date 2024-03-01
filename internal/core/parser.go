package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func parseAssignmentExpr(expr string) (string, string, error) {
	// check if assingment exist
	index := strings.Index(expr, "=")
	if index == -1 {
		return "", "", errors.New("Invalid syntax for `var` tag, expected variable assignment")
	}
	// split into lhs and rhs
	lhs := expr[:index]
	rhs := expr[index+1:]
	// check if expression is empty
	if len(lhs) == 0 || len(rhs) == 0 {
		return "", "", errors.New("The LHS and RHS of expression must not be empty")
	}
	return lhs, rhs, nil
}

func parseTemplateExpr(tokens []string) (*block, error) {
	var data []string
	// check no template name given
	if len(tokens) < 2 {
		return nil, errors.New("Invalid syntax for `template` tag, no template name given")
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
}

func parseExpandExpr(tokens []string) (*block, error) {
	var data []string
	// check no template name given
	if len(tokens) < 2 {
		return nil, errors.New("Invalid syntax for `expand` tag, no template name given")
	}
	data = append(data, tokens[1])
	// check not correct form of expression
	if len(tokens) >= 3 {
		for _, assignExpr := range tokens[2:] {
			key, val, err := parseAssignmentExpr(assignExpr)
			if err != nil {
				return nil, errors.New("Invalid syntax for `expand` tag, invalid assignment expression")
			}
			data = append(data, key)
			data = append(data, val)
		}
	}
	return &block{
		blockType: blockTypeExapand,
		data:      data,
	}, nil
}

func parseContentExpr(tokens []string) (*block, error) {
	return &block{
		blockType: blockTypeContent,
	}, nil
}

func parseUseExpr(tokens []string) (*block, error) {
	var data []string
	// check invalid length of variables
	if len(tokens) != 2 {
		return nil, errors.New("Invalid syntax for `use` tag, no variable name given")
	}
	data = append(data, tokens[1])
	return &block{
		blockType: blockTypeUse,
		data:      data,
	}, nil
}

func parseForExpr(tokens []string) (*block, error) {
	var data []string
	// check invalid length
	if len(tokens) != 4 {
		return nil, errors.New("Invalid syntax for `for` tag")
	}
	// check for in
	if tokens[2] != "in" {
		return nil, errors.New("Invalid syntax for `for` tag, keyword `in` was expected")
	}
	data = append(data, tokens[1])
	data = append(data, tokens[3])
	return &block{
		blockType: blockTypeFor,
		data:      data,
	}, nil
}

func parseEndForExpr(tokens []string) (*block, error) {
	return &block{
		blockType: blockTypeEndFor,
	}, nil
}

func parseVarExpr(tokens []string) (*block, error) {
	var data []string
	// check invalid length
	if len(tokens) != 2 {
		return nil, errors.New("Invalid syntax for `var` tag")
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
}

func parseOutOnlyExpr(tokens []string) (*block, error) {
	return &block{
		blockType: blockTypeOutOnly,
	}, nil
}

// Parses non HTML blocks between start and end
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
		return parseTemplateExpr(tokens)
	case "expand":
		return parseExpandExpr(tokens)
	case "content":
		return parseContentExpr(tokens)
	case "use":
		return parseUseExpr(tokens)
	case "for":
		return parseForExpr(tokens)
	case "endfor":
		return parseEndForExpr(tokens)
	case "var":
		return parseVarExpr(tokens)
	case "outonly":
		return parseOutOnlyExpr(tokens)
	default:
		return nil, errors.New(fmt.Sprintf("Unrecognized block type '%s'", blockTypeStr))
	}
}

// Convert HTML data to chain of blocks
func parseHtmlBlockChain(data *[]byte) (*blockChain, error) {
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
