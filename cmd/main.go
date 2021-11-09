package main

import (
	"fmt"
	"github.com/ManizhaM/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil{
		fmt.Println(account, err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	if err != nil {
		fmt.Println(account, err)
		return
	}

	fmt.Println(account.ID)//10
}