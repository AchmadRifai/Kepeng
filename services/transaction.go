package services

import "main/galats"

func NewTransaction() TransactionRouter {
	return TransactionRouter{}
}

type TransactionRouter struct {
	galats.PesimisticLocking
}
