-- +goose Up
-- +goose StatementBegin
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_item (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION add_item_to_cart(
    p_cart_id INT,
    product_name VARCHAR,
    price DECIMAL
)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    distinct_count INT;
    new_id INT;
BEGIN
    IF NOT EXISTS(SELECT 1 FROM carts WHERE id = p_cart_id) THEN
       RAISE EXCEPTION 'Cart % does not exist', p_cart_id;
    END IF;

    IF product_name IS NULL OR TRIM(product_name) = '' THEN
       RAISE EXCEPTION 'product name cannot be blank';
    END IF;

    IF price <= 0 OR price IS NULL THEN
       RAISE EXCEPTION 'incorrect price information: %', price;
    END IF;

    SELECT COUNT(DISTINCT product) INTO distinct_count
    FROM cart_item
    WHERE cart_id = p_cart_id;

    IF distinct_count >= 5 AND NOT EXISTS (
       SELECT 1 FROM cart_item
       WHERE cart_id = p_cart_id AND product = product_name) THEN
    RAISE EXCEPTION 'cart cannot consist more than 5 distinct products';
    END IF;

    INSERT INTO cart_item(cart_id, product, price)
    VALUES (p_cart_id, product_name, price)
    RETURNING id INTO new_id;

    RETURN new_id;
END;
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE carts, cart_item;
-- +goose StatementEnd
