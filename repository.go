package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func CreateConnectionPool(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(25)

	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Conexão com o banco de dados PostgreSQL estabelecida com sucesso!")
	return nil
}

func GetClientById(id int) (Cliente, error) {
	var cliente Cliente
	rows, err := db.Query("SELECT * FROM clientes WHERE id = $1", id)
	if err != nil {
		return cliente, err
	}

	for rows.Next() {
		err := rows.Scan(&cliente.Id, &cliente.Limite, &cliente.SaldoInicial)
		if err != nil {
			return cliente, err
		}
		return cliente, nil
	}

	err = rows.Err()
	if err != nil {
		return cliente, err
	}

	defer rows.Close()
	return cliente, fmt.Errorf("Cliente com ID %d não encontrado", id)
}

func SaveCliente(cliente Cliente) (bool, error) {
	_, err := db.Query("UPDATE clientes SET saldo_inicial = $1 WHERE id = $2", cliente.SaldoInicial, cliente.Id)
	if err != nil {
		return false, errors.New(err.Error())
	}
	return true, nil

}

func SaveTransaction(id int, transacaoDTO TransacaoDTO, cliente Cliente) (TransacaoRespostaDTO, error) {

	if transacaoDTO.Tipo != "c" && transacaoDTO.Tipo != "d" {
		return TransacaoRespostaDTO{}, errors.New("Tipo inválido")
	}

	if len(transacaoDTO.Descricao) < 1 || len(transacaoDTO.Descricao) > 10 {
		return TransacaoRespostaDTO{}, errors.New("Descrição deve apenas conter entre 1 a 10 caracteres")
	}

	novoLimite := cliente.Limite - -cliente.SaldoInicial
	if transacaoDTO.Valor > novoLimite {
		return TransacaoRespostaDTO{}, errors.New("Valor da transação excede o limite")
	}

	_, err := db.Query("INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)",
		id,
		transacaoDTO.Valor,
		transacaoDTO.Tipo,
		transacaoDTO.Descricao,
		time.Now())
	if err != nil {
		return TransacaoRespostaDTO{}, err
	}

	var respostaDTO TransacaoRespostaDTO
	respostaDTO.Limite = cliente.Limite
	respostaDTO.Saldo = cliente.SaldoInicial

	SaveCliente(cliente)

	return respostaDTO, nil

}

func GetExtratoByClienteId(id int) (ExtratoFinalRespostaDTO, error) {
	rows, err := db.Query("SELECT limite, saldo_inicial, valor, descricao, tipo, realizada_em FROM clientes c LEFT JOIN transacoes t ON t.cliente_id = c.id WHERE c.id = $1 ORDER BY t.id DESC LIMIT 10", id)
	if err != nil {
		return ExtratoFinalRespostaDTO{}, err
	}

	if id > 5 {
		return ExtratoFinalRespostaDTO{}, errors.New("Cliente não encontrado")
	}

	res := ExtratoFinalRespostaDTO{
		UltimasTransacoes: make([]ExtratoRespostaDTO, 0, 10),
	}

	for rows.Next() {
		var carteira CarteiraRespostaDTO
		var transacao ExtratoRespostaDTO
		err = rows.Scan(&carteira.Limite, &carteira.Saldo, &transacao.Valor, &transacao.Descricao, &transacao.Tipo, &transacao.RealizadaEm)
		if err != nil {
			if carteira.Limite != 0 {
				res.Saldo.Saldo = carteira.Saldo
				res.Saldo.Limite = carteira.Limite
				res.Saldo.DataExtrato = time.Now()
			}

			fmt.Println(fmt.Errorf("Unable to scan row %v", err))
		}

		res.Saldo.Saldo = carteira.Saldo
		res.Saldo.Limite = carteira.Limite
		res.Saldo.DataExtrato = time.Now()
		res.UltimasTransacoes = append(res.UltimasTransacoes, transacao)
	}

	// Close the rows after iterating over it
	defer rows.Close()
	return res, nil
}
