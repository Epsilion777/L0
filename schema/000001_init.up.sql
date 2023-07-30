CREATE TABLE orders (
	order_uid UUID PRIMARY KEY,
	track_number VARCHAR(100),
	entry VARCHAR(100),
	locale VARCHAR(100),
	internal_signature VARCHAR(100),
	customer_id VARCHAR(100),
	delivery_service VARCHAR(100),
	shardkey VARCHAR(100),
	sm_id INTEGER,
	date_created VARCHAR(100),
	oof_shard VARCHAR(100)
);

CREATE TABLE deliveries (
    order_uid UUID REFERENCES orders (order_uid) ON DELETE CASCADE,
    name VARCHAR(100),
    phone VARCHAR(100),
    zip VARCHAR(100),
    city VARCHAR(100),
    address VARCHAR(100),
    region VARCHAR(100),
    email VARCHAR(100)
);

CREATE TABLE payments (
    order_uid UUID REFERENCES orders (order_uid) ON DELETE CASCADE,
    transaction VARCHAR(100),
    request_id VARCHAR(100),
    currency VARCHAR(100),
    provider VARCHAR(100),
    amount INTEGER,
    payment_dt INTEGER,
    bank VARCHAR(100),
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE items (
    order_uid UUID REFERENCES orders (order_uid) ON DELETE CASCADE,
    chrt_id INTEGER,
    track_number VARCHAR(100),
    price INTEGER,
    rid VARCHAR(100),
    name VARCHAR(100),
    sale INTEGER,
    size VARCHAR(100),
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR(100),
    status INTEGER
);
