package main

type NewAccount struct {
	Accountname   string
	Currency      string
	Initialamount int
}

type HistoryElement struct {
	Name     string
	Id       int
	Category string
	Price    int
	Time     string
}
