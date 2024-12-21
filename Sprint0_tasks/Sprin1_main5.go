package main

type account struct {
	balance int
	owner   string

	SetBalance int // перересоздать боланс
	Deposit    int // + деньги
	Withdraw   int // - деньги
	GetBalance int // узнать боланс
}

func NewAccount(balance int, owner string) *account {

}
