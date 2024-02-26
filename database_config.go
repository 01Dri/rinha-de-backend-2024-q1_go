package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DbConfig struct {
	Name     string
	Port     int
	User     string
	Password string
	Url      string
}

func startConnection(config DbConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.User, config.Password, config.Name, config.Port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Conexão com o banco de dados PostgreSQL estabelecida com sucesso!")
	return db, nil
}

func getClientById(id int, conn *sql.DB) (Cliente, error) {
	var cliente Cliente

	rows, err := conn.Query("SELECT * FROM clientes WHERE id = $1", id)
	if err != nil {
		return cliente, err
	}

	defer rows.Close()

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

	return cliente, fmt.Errorf("Cliente com ID %d não encontrado", id)
}

func saveTransaction(id int, conn *sql.DB, transacaoDTO TransacaoDTO) (bool, error) {
	_, err := conn.Query("INSERT INTO transacoes (client_id, valor, tipo, descricao, realizada_em) VALUES ($1, $2, $3, $4, $5)",
		id,
		transacaoDTO.Valor,
		transacaoDTO.Tipo,
		transacaoDTO.Descricao,
		time.Now())
	if err != nil {
		return false, err
	}
	return true, nil

}

func getExtratoByClienteId(id int, conn *sql.DB) (ExtratoFinalRespostaDTO, error) {
	rows, err := conn.Query("SELECT limite, saldo_inicial, valor, descricao, tipo, realizada_em FROM clientes c LEFT JOIN transacoes t ON t.cliente_id = c.id WHERE c.id = $1	 ORDER BY t.id DESC LIMIT 10", id)
	if err != nil {
		return ExtratoFinalRespostaDTO{}, err
	}
	res := ExtratoFinalRespostaDTO{
		UltimasTransacoes: make([]ExtratoRespostaDTO, 0, 10),
	}

	for rows.Next() {
		var carteira CarteiraRespostaDTO
		var transacao ExtratoRespostaDTO
		err = rows.Scan(&carteira.Limite, &carteira.Saldo, &transacao.Saldo, &transacao.Descricao, &transacao.Tipo, &transacao.RealizadaEm)
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
	return res, nil

}
