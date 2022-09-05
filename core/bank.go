package bank

import (
	"errors"
	"fmt"
)

type Customer struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}
type Account struct {
	Customer
	Number  int32   `json:"number"`
	Balance float64 `json:"balance"`
}

type Bank interface {
	State() string
}

//func CreateAccount(name, address, phone string, number int32) *Account {
//	return &Account{
//		Customer: Customer{
//			Name:    name,
//			Address: address,
//			Phone:   phone,
//		},
//		Number:  number,
//		Balance: 0,
//	}
//}

func (a *Account) State() string {
	return fmt.Sprintf("%v - %v - %v", a.Number, a.Name, a.Balance)
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("the amount to deposit should be greater than zero")
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("the amount to deposit should be greater than zero")
	}

	if amount > a.Balance {
		return errors.New("insufficient balance")
	}
	a.Balance -= amount
	return nil
}

func (src *Account) Transfer(dest *Account, amount float64) error {
	if amount <= 0 {
		return errors.New("amount cannot be less than zero")
	}

	if src.Balance < amount {
		return fmt.Errorf("%s balance is insufficient", src.Name)
	}
	err := src.Withdraw(amount)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = dest.Deposit(amount)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func Statement(b Bank) string {
	return b.State()
}
