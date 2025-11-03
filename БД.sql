-- Role: order_service_user

CREATE ROLE order_service_user WITH
  LOGIN
  SUPERUSER
  INHERIT
  CREATEDB
  CREATEROLE
  NOREPLICATION
  BYPASSRLS
  ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:eQs9UiaM2R8JkvHGtYK9Mw==$X+t+5yxHHOSghV0s8bBpIs14j5XU82APF1+Cagkwjn4=:v0g8nFCKW8HsNYEq8zU+C5iRy9Hbg1mPC92rij5S+8w=';


-- Создание таблиц
CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(50),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(100),
    shardkey VARCHAR(50),
    sm_id INTEGER,
    date_created TIMESTAMP WITH TIME ZONE,
    oof_shard VARCHAR(50)
);

CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR(255),
    phone VARCHAR(50),
    zip VARCHAR(50),
    city VARCHAR(255),
    address TEXT,
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE TABLE payment (
    transaction VARCHAR(255) PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(100),
    amount INTEGER,
    payment_dt BIGINT,
    bank VARCHAR(100),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT,
    track_number VARCHAR(255),
    price INTEGER,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INTEGER,
    size VARCHAR(50),
    total_price INTEGER,
    nm_id BIGINT,
    brand VARCHAR(255),
    status INTEGER
);

-- Создание индексов
CREATE INDEX idx_orders_order_uid ON orders(order_uid);
CREATE INDEX idx_delivery_order_uid ON delivery(order_uid);
CREATE INDEX idx_payment_order_uid ON payment(order_uid);
CREATE INDEX idx_items_order_uid ON items(order_uid);