package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var (
	GenesisBlock     = &Block{0, time.Now().String(), 0, "", "", ""}
	MasterBlockchain = []*Block{GenesisBlock}
	mclock           sync.RWMutex
)

type Block struct {
	// Index is the position of the data record in the blockchain
	Index int

	// Timestamp is automatically determined and is the time the data is written
	Timestamp string

	// BPM or beats per minute, is your pulse rate
	BPM int

	// Hash is a SHA256 identifier representing this data record
	Hash string

	// PrevHash is the SHA256 identifier of the previous record in the chain
	PrevHash string

	// Information record in this block
	Content string
}

func CalculateHash(block *Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateBlock(oldBlock *Block, BPM int) (*Block, error) {
	newBlock := &Block{}

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)

	return newBlock, nil
}

func GenerateBlockWithContent(oldBlock *Block, BPM int, content string) (*Block, error) {
	newBlock := &Block{}

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)
	newBlock.Content = content

	return newBlock, nil
}

func IsBlockValid(newBlock, oldBlock *Block) bool {
	//time.Sleep(time.Second)

	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func (n *Node) ReplaceChain(newBlocks []*Block) {
	if len(newBlocks) > len(n.chain) {
		n.chain = newBlocks
	}
}

func PrintBlockChain() {
	spew.Dump(MasterBlockchain)
}
