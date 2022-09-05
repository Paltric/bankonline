package main

import (
	"bank"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var accounts = map[float64]*bank.Account{}

func main() {
	accounts[1001] = &bank.Account{
		Customer: bank.Customer{
			Name:    "mike",
			Address: "London",
			Phone:   "028-666123",
		},
		Number:  1001,
		Balance: 0,
	}
	http.HandleFunc("/statement", statement)
	http.HandleFunc("/deposit", deposit)
	http.HandleFunc("/withdraw", withdraw)
	http.HandleFunc("/transfer", transfer)
	log.Fatalln(http.ListenAndServe(":8000", nil))
}

func transfer(w http.ResponseWriter, req *http.Request) {
	numberqs := req.URL.Query().Get("number")
	amountqs := req.URL.Query().Get("amount")
	destqs := req.URL.Query().Get("dest")

	if numberqs == "" {
		fmt.Fprintf(w, "account number is missing!")
		return
	}

	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
		return
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
	} else if dest, err := strconv.ParseFloat(destqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account destination number!")
	} else {
		if accountA, ok := accounts[number]; !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else if accountB, ok := accounts[dest]; !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", dest)
		} else {
			err := accountA.Transfer(accountB, amount)
			if err != nil {
				fmt.Fprintf(w, "%v", err)
			} else {
				fmt.Fprintf(w, accountA.State())
			}
		}
	}
}

func withdraw(w http.ResponseWriter, req *http.Request) {
	numberqs := req.URL.Query().Get("number")
	amountqs := req.URL.Query().Get("amount")

	if numberqs == "" {
		fmt.Fprintf(w, "account number is missing!")
		return
	}

	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
		return
	} else if amount, err := strconv.ParseFloat(amountqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
	} else {
		if account, ok := accounts[number]; !ok {
			fmt.Fprintf(w, "Account with number %v can't be found!", number)
		} else {
			err := account.Withdraw(amount)
			if err != nil {
				fmt.Fprintf(w, "A error found: %v", err)
			} else {
				fmt.Fprintf(w, account.State())
			}
		}
	}
}

func deposit(w http.ResponseWriter, req *http.Request) {
	numberqs := req.URL.Query().Get("number")
	depositqs := req.URL.Query().Get("amount")

	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
		return
	} else if amount, err := strconv.ParseFloat(depositqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid amount number!")
		return
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "account number is not found!")
		} else {
			if err := account.Deposit(amount); err != nil {
				fmt.Fprintf(w, "account deposit with error: %v", err)
			} else {
				fmt.Fprintf(w, account.State())
			}
		}
	}
}

func statement(w http.ResponseWriter, req *http.Request) {
	numberqs := req.URL.Query().Get("number")

	if number, err := strconv.ParseFloat(numberqs, 64); err != nil {
		fmt.Fprintf(w, "Invalid account number!")
		return
	} else {
		account, ok := accounts[number]
		if !ok {
			fmt.Fprintf(w, "account with number %v can't be found!", number)
		} else {
			c := &CustomAccount{account}
			json.NewEncoder(w).Encode(bank.Statement(c))
		}
	}
}

type CustomAccount struct {
	*bank.Account
}

func (c *CustomAccount) State() string {
	ac, err := json.Marshal(c)
	if err != nil {
		return err.Error()
	}

	return string(ac)
}
