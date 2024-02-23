package main

import (
	"database/sql"
	"fmt"

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
