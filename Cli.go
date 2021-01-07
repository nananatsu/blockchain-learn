package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI 命令行工具
type CLI struct {
}

func (cli *CLI) createBlockchain(address string) {

	bc := CreateBlockChain(address)
	bc.db.Close()

	UTXOSet := UTXOSet{bc}
	UTXOSet.Reindex()

	fmt.Println("创建链成功！")
}

func (cli *CLI) getBalance(address string) {

	if !ValidateAddress(address) {
		log.Panic("非法的地址！")
	}

	bc := NewBlockchain()
	defer bc.db.Close()

	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]

	UTXOSet := UTXOSet{bc}

	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("余额: %s:%d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {

	if !ValidateAddress(from) {
		log.Panic("发送者地址非法")
	}
	if !ValidateAddress(to) {
		log.Panic("接收者地址非法")
	}

	bc := NewBlockchain()
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	newBlock := bc.MineBlock([]*Transaction{tx})

	UTXOSet := UTXOSet{bc}
	UTXOSet.Update(newBlock)

	fmt.Println("交易成功")
}

func (cli *CLI) printChain() {
	bc := NewBlockchain()
	defer bc.db.Close()

	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ 块 %x ============\n", block.Hash)
		fmt.Printf("前一块: %x\n", block.PrevBlockHash)
		pow := NewProofOfWork(block)
		fmt.Printf("工作量证明: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}

	}
}

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("钱包地址: %s\n", address)
}

func (cli *CLI) listAddresses() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.GetAddress()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
}

// Run 运行命令行工具
func (cli *CLI) Run() {

	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listAddresses", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "要查询余额的地址")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "发送挖矿奖励的地址")
	sendFrom := sendCmd.String("from", "", "源钱包地址")
	sendTo := sendCmd.String("to", "", "目的钱包地址")
	sendAmount := sendCmd.Int("amount", 0, "交易金额")

	switch os.Args[1] {
	case "printChain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createBlockChain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listAddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		os.Exit(1)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}

}
