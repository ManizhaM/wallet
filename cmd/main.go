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

	err = svc.Deposit(account.ID, 2000)
	if err != nil {
		fmt.Println(account, err)
		return
	}

	pay, _ := svc.Pay(1, 1000,"bank")
	//err2 := svc.Reject(pay.ID)
	
	payment2,err := svc.Repeat(pay.ID)
	payment, err := svc.FindAccountByID(payment2.AccountID)
	//fmt.Println(err2)//10
	fmt.Println(payment)//10
}