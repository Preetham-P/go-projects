package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ISBN        string `json:"isbn"`
	Author      string `json:"author"`
	PublishDate string `json:"publishdate"`
}

type Block struct {
	Pos       int
	Data      BookCheckout
	TimeStamp string
	PrevHash  string
	Hash      string
}

type BookCheckout struct {
	BookID       string `json:"bookid"`
	IsGenesis    bool   `json:"isgenesis"`
	User         string `json:"user"`
	CheckOutDate string `json:"checkoutdate"`
}

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prevBlock, data)
	bc.blocks = append(bc.blocks, block)
}

func (b *Block) GenerateHash() {
	bytes, _ := json.Marshal(b.Data)

	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()

	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))

}

func CreateBlock(prevBlock *Block, data BookCheckout) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.PrevHash = prevBlock.Hash
	block.TimeStamp = time.Now().String()
	block.GenerateHash()
	return block
}

var blockChain *BlockChain

func getBlockChain(w http.ResponseWriter, r *http.Request) {
	bytedata, err := json.MarshalIndent(blockChain.blocks, "", " ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(bytedata))
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{GenesisBlock()}}
}

func writeBlock(w http.ResponseWriter, r *http.Request) {

	var bookCheckout BookCheckout

	if err := json.NewDecoder(r.Body).Decode(&bookCheckout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not add checkout block"))
		return
	}

	blockChain.AddBlock(bookCheckout)

}

func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf(err.Error())
		w.Write([]byte("Could not create the book"))
		return
	}

	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(book, "", " ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not create a book"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func main() {
	blockChain = NewBlockChain()

	fmt.Println("This is simple bock chain demo")

	routes := mux.NewRouter()
	routes.HandleFunc("/", getBlockChain).Methods("GET")
	routes.HandleFunc("/", writeBlock).Methods("POST")
	routes.HandleFunc("/new", newBook).Methods("POST")

	go func() {

		for _, block := range blockChain.blocks {
			fmt.Printf("Prev. hash: %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data: %v\n", string(bytes))
			fmt.Printf("Hash: %x\n", block.Hash)
			fmt.Println()
		}

	}()

	log.Println("starting the server on 9090")

	log.Fatal(http.ListenAndServe(":9090", routes))
}
