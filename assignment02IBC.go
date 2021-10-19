//package main
package assignment02IBC

import (
	"crypto/sha256"
	"fmt"
)

const miningReward = 100
const rootUser = "Satoshi"

type BlockData struct {
	Title    string
	Sender   string
	Receiver string
	Amount   int
}
type Block struct {
	Data        []BlockData
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

func CalculateBalance(userName string, chainHead *Block) int {
	var balance int = 0
	for trav := chainHead; trav != nil; trav = trav.PrevPointer {
		for i := 0; i < len(trav.Data); i++ {
			if trav.Data[i].Sender == userName {
				balance = balance - trav.Data[i].Amount
			}
			if trav.Data[i].Receiver == userName {
				balance = balance + trav.Data[i].Amount
			}
		}
	}
	return balance
}

func CalculateHash(inputBlock *Block) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v", inputBlock.Data))))
}

func VerifyTransaction(transaction *BlockData, chainHead *Block) bool {
	//check if sender has the balance to make that transaction
	canGive := CalculateBalance(transaction.Sender, chainHead)
	if canGive >= transaction.Amount {
		return true
	}
	return false
}

func InsertBlock(blockData []BlockData, chainHead *Block) *Block {
	// for every transaction Satoshi gets 100 Coinbased balance
	//catering for that
	blockData = append(blockData, []BlockData{{Title: "Coinbased", Sender: "System", Receiver: rootUser, Amount: miningReward}}...)

	var newBlock *Block = new(Block)
	newBlock.Data = blockData
	if chainHead != nil {
		newBlock.PrevHash = chainHead.CurrentHash
		newBlock.PrevPointer = chainHead
	} else {
		newBlock.PrevHash = ("0")
		newBlock.PrevPointer = nil
	}
	newBlock.CurrentHash = CalculateHash(newBlock)
	for i := 0; i < len(blockData); i++ {
		// If any transaction has less balance than amount sent unless the sender is "System", transaction is invalid
		if !VerifyTransaction(&blockData[i], chainHead) && blockData[i].Sender != "System" {
			if blockData[i].Sender == rootUser {
				fmt.Println("ERROR:", blockData[i].Sender, " has ", CalculateBalance(blockData[i].Sender, chainHead), " coins - ", CalculateBalance(blockData[i].Sender, chainHead)-CalculateBalance(blockData[i].Sender, newBlock)+miningReward, " were needed!")
			} else {
				fmt.Println("ERROR:", blockData[i].Sender, " has ", CalculateBalance(blockData[i].Sender, chainHead), " coins - ", CalculateBalance(blockData[i].Sender, chainHead)-CalculateBalance(blockData[i].Sender, newBlock), " were needed!")
			}
			return chainHead
		}
	}
	// FOR BONUS, traversing the whole transaction chain to check everyone's balances after adding the new block
	// If anyone has insufficient balance, block will not be added
	var invalid bool = false
	var repeat bool = false
	var doubleSpenders []string
	for i := 0; i < len(blockData); i++ {
		// If any transaction has less balance than amount sent unless the sender is "System", transaction is invalid
		if CalculateBalance(blockData[i].Sender, newBlock) < 0 && blockData[i].Sender != "System" {
			for j := 0; j < len(doubleSpenders); j++ {
				if doubleSpenders[j] == blockData[i].Sender {
					repeat = true
				}
			}
			if repeat == false {
				fmt.Println()
				//fmt.Println("------------->>>BONUS WORK<<<-------------")
				fmt.Println("ERROR:", blockData[i].Sender, " has insufficient balance to carry out the transactions")
				//fmt.Println(blockData[i].Sender, "is double spending!")
				doubleSpenders = append(doubleSpenders, blockData[i].Sender)
				invalid = true
			}
		}
		repeat = false
	}

	if invalid == true {
		fmt.Println()
		fmt.Println("New block not added to the chain since transactions were invalid.")
		fmt.Println()
		//Return the chain without the invalid new block
		return chainHead
	}

	return newBlock
}

func ListBlocks(chainHead *Block) {
	for trav := chainHead; trav != nil; trav = trav.PrevPointer {

		for i := 0; i < len(trav.Data); i++ {
			fmt.Println("Title:", trav.Data[i].Title, "Sender:", trav.Data[i].Sender, "Receiver:", trav.Data[i].Receiver, "Amount:", trav.Data[i].Amount)
		}
		fmt.Println()
	}
}

func VerifyChain(chainHead *Block) {
	prev := chainHead.PrevPointer

	for trav := chainHead; trav.PrevPointer != nil; trav = trav.PrevPointer {
		if prev.CurrentHash != trav.PrevHash {
			fmt.Println("Chain Compromised!")
			break
		}
		prev = chainHead.PrevPointer
	}
}

func PremineChain(chainHead *Block, numBlocks int) *Block {
	//var preminedChain *Block
	for i := 0; i < numBlocks; i++ {
		premined := []BlockData{{Title: "Premined", Sender: "nil", Receiver: "nil", Amount: 0}}
		chainHead = InsertBlock(premined, chainHead)
	}
	//return preminedChain
	return chainHead
}

/*
func main() {
	var chainHead *Block
	//This insertion is invalid as Alice is neither miner nor has enough coins for the transaction, pay 50 from Alice to Bob
	aliceToBob := []BlockData{{Title: "ALice2Bob", Sender: "Alice", Receiver: "Bob", Amount: 50}}
	chainHead = InsertBlock(aliceToBob, chainHead)

	//Lets mine some blocks to start the chain and check Satoshi's balance
	chainHead = PremineChain(chainHead, 2)
	//fmt.Printf("Satoshi's balance: %v\n", a2.CalculateBalance("Satoshi", chainHead))

	//Now Satoshi can send some coins to Alice
	SatoshiToAlice := []BlockData{{Title: "SatoshiToAlice", Sender: "Satoshi", Receiver: "Alice", Amount: 50}}
	chainHead = InsertBlock(SatoshiToAlice, chainHead)

	//We can verify this by checking balances once again and listing the chain
	/*
	   fmt.Printf("Satoshi's balance: %v\n", a2.CalculateBalance("Satoshi", chainHead))
	   fmt.Printf("Alice's balance: %v\n", a2.CalculateBalance("Alice", chainHead))
	   a2.ListBlocks(chainHead)
*/
//Alice can then make the transactions using her coins, She can make multiple
//transactions at once, notice that field Data has type []BlockData in Block Struct
/*AliceToBobCharlie := []BlockData{{Title: "ALice2Bob", Sender: "Alice", Receiver: "Bob", Amount: 20}, {Title: "ALice2Charlie", Sender: "Alice", Receiver: "Charlie", Amount: 10}}
	chainHead = InsertBlock(AliceToBobCharlie, chainHead)

	//We can verify this by checking balances once again and listing the chain
	ListBlocks(chainHead)

	fmt.Printf("Satoshi's balance: %v ", CalculateBalance("Satoshi", chainHead))
	fmt.Printf("Alice's balance: %v ", CalculateBalance("Alice", chainHead))
	fmt.Printf("Charlie's balance: %v\n", CalculateBalance("Charlie", chainHead))

	//Finally the transaction verification fails if any of the transaction is invalid
	oneInvalidoneValid := []BlockData{{Title: "ALice2EZ", Sender: "Alice", Receiver: "Bob", Amount: 100}, {Title: "Satoshi2EZ", Sender: "Satoshi", Receiver: "EZ", Amount: 200}}
	chainHead = InsertBlock(oneInvalidoneValid, chainHead)

	//Bonus (2 absolutes) - Fix the erroneous behavior below
	//The transactions are valid individually but when applied to chain Alice's balance
	//become negative :(
	bonusTransactions := []BlockData{{Title: "ALice2Bob", Sender: "Alice", Receiver: "Bob", Amount: 15}, {Title: "AliceToEZ", Sender: "Alice", Receiver: "EZ", Amount: 15}}
	chainHead = InsertBlock(bonusTransactions, chainHead)
	fmt.Printf("Alice's balance: %v\n", CalculateBalance("Alice", chainHead))

}*/
