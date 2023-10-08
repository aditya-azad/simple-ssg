package core

//import "github.com/aditya-azad/simple-ssg/pkg/logging"

const (
	BLOCK_CONTENT = iota
	BLOCK_USE
	BLOCK_FOR
	BLOCK_OUT_ONLY
	BLOCK_EXPAND
	BLOCK_TEMPLATE
	BLOCK_SENTINEL
)

type Block struct {
	Type int
	Next *Block
	Prev *Block
}

type BlockChain struct {
	sentinel *Block
}

func New() *BlockChain {
	// Create and return a new block chain
	sentinel := Block{}
	sentinel.Type = BLOCK_SENTINEL
	sentinel.Next = &sentinel
	sentinel.Prev = &sentinel
	return &BlockChain{&sentinel}
}

func (bc *BlockChain) Append(b *Block) {
	// Insert new block at the end of the list
	last := bc.sentinel.Prev
	b.Next = bc.sentinel
	b.Prev = last
	last.Next = b
	bc.sentinel.Prev = b
}

func (bc *BlockChain) AppendLeft(b *Block) {
	// Insert new block at the head of the list
	next := bc.sentinel.Next
	b.Next = next
	b.Prev = bc.sentinel
	next.Prev = b
	bc.sentinel.Next = b
}

func (bc *BlockChain) PopLeft() *Block {
	// Remove a block from the start of list, return error if not present
}

func (bc *BlockChain) Pop() *Block {
	// Remove a block from the end of list, return error if not present
}

func (bc *BlockChain) Eject() (*Block, *Block) {
	// Remove sentinel and return the head and tail of the list, return error if list is empty
}
