-- +goose Up
-- +goose StatementBegin
CREATE TABLE carts (
    id SERIAL PRIMARY KEY
);

CREATE TABLE cart_item (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE carts, cart_item;
-- +goose StatementEnd
