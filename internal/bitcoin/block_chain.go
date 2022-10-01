package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// ref: https://mycoralhealth.medium.com/code-your-own-blockchain-in-less-than-200-lines-of-go-e296282bcffc
var (
	genesisBlock     = &Block{0, time.Now().String(), 0, "", ""}
	MasterBlockchain = []*Block{genesisBlock}
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
}

func calculateHash(block *Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock *Block, BPM int) (*Block, error) {
	var newBlock *Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}

func isBlockValid(newBlock, oldBlock *Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []*Block) {
	if len(newBlocks) > len(MasterBlockchain) {
		MasterBlockchain = newBlocks
	}
}

func PrintBlockChain() {
	spew.Dump(MasterBlockchain)
}
