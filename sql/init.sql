CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Tabela de Usuários (Simplificada para suportar a FK)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
);

-- 2. Tabela de Eventos
-- O bloqueio (FOR UPDATE) será feito nesta tabela
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    total_capacity INT NOT NULL,
    base_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Tabela de Pedidos (Orders)
CREATE TABLE orders (
    id UUID PRIMARY KEY, -- Gerado pelo Go
    event_id UUID NOT NULL REFERENCES events(id),
    user_id UUID NOT NULL REFERENCES users(id),
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL, -- PENDING, PAID, CANCELED
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Tabela de Ingressos (Tickets)
CREATE TABLE tickets (
    id UUID PRIMARY KEY, -- Gerado pelo Go
    event_id UUID NOT NULL REFERENCES events(id),
    order_id UUID NOT NULL REFERENCES orders(id),
    price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL
);

-- Índices para performance da Query de Contagem
CREATE INDEX idx_tickets_event_id ON tickets(event_id);
CREATE INDEX idx_orders_status ON orders(status);

-- === SEED DATA (Dados iniciais para teste) ===

-- Cria um usuário de teste
INSERT INTO users (id, name, email) 
VALUES ('a1b2c3d4-e5f6-7890-1234-567890abcdef', 'Tester User', 'tester@example.com');

-- Cria um evento "Rock in Go" com capacidade de 10 pessoas
INSERT INTO events (id, name, total_capacity, base_price) 
VALUES ('b2c3d4e5-f678-9012-3456-7890abcdef12', 'Rock in Go 2026', 10, 100.00);