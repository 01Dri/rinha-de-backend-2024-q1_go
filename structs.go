package main

import (
	"time"
)

type Cliente struct {
	Id           int `json:"id"`
	Limite       int `json:"limite"`
	SaldoInicial int `json:"saldo_inicial"`
}

type TransacaoDTO struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

type TransacaoRespostaDTO struct {
	Limite int `json:"limite"`
	Saldo  int `json:"saldo"`
}

type CarteiraRespostaDTO struct {
	Saldo       int       `json:"saldo"`
	Limite      int       `json:"limite"`
	DataExtrato time.Time `json:"data_extrato"`
}

type ExtratoRespostaDTO struct {
	Valor       int       `json:"valor"`
	Tipo        string    `json:"tipo"`
	Descricao   string    `json:"descricao"`
	RealizadaEm time.Time `json:"realizada_em"`
}

type ExtratoFinalRespostaDTO struct {
	Saldo             CarteiraRespostaDTO  `json:"saldo"`
	UltimasTransacoes []ExtratoRespostaDTO `json:"ultimas_transacoes"`
}
