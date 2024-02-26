CREATE DATABASE rinha_back_end;

\c rinha_back_end;

CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    limite BIGINT,
    saldo_inicial BIGINT
);

CREATE TABLE transacoes (
    id SERIAL PRIMARY KEY,
    cliente_id INT,
    valor BIGINT,
    tipo CHAR(1),
    descricao TEXT,
    realizada_em TIMESTAMP,
    FOREIGN KEY (cliente_id) REFERENCES clientes(id)
);

INSERT INTO clientes (limite, saldo_inicial) VALUES
    (100000, 0),
    (80000, 0),
    (1000000, 0),
    (10000000, 0),
    (500000, 0);

-- Inserir dados na tabela transacoes
INSERT INTO transacoes (cliente_id, valor, tipo, descricao, realizada_em) VALUES
    (1, 10, 'c', 'descricao', '2024-01-17T02:34:38.543030Z'),
    (1, 90000, 'd', 'descricao', '2024-01-17T02:34:38.543030Z'); id |  limite  | saldo_inicial 
