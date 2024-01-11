package models

type Account struct {
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"desc"`
}

type Accounts struct {
	Datas []Account `json:"datas"`
}

type Asset struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Increase float64 `json:"increase"`
	Decrease float64 `json:"decrease"`
}

type Assets struct {
	Datas []Asset `json:"datas"`
}

type Liability struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Interest float64 `json:"interest"`
	DueDate  string  `json:"dueDate"`
}

type Liabilities struct {
	Datas []Liability `json:"datas"`
}

type Income struct {
	Name    string `json:"name"`
	Created string `json:"created"`
}

type Incomes struct {
	Datas []Income `json:"datas"`
}

type Outcome struct {
	Name    string `json:"name"`
	Created string `json:"created"`
}

type Outcomes struct {
	Datas []Outcome `json:"datas"`
}

type Transaction struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Time        string  `json:"time"`
}

type Transactions struct {
	Datas []Transaction `json:"datas"`
}

type TypeAccount struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TypeAccounts struct {
	Datas []TypeAccount `json:"datas"`
}
