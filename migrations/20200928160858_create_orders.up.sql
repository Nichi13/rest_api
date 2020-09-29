CREATE TABLE orders (
                        id BIGSERIAL NOT NULL PRIMARY KEY,
                        status VARCHAR NOT NULL,
                        number BIGSERIAL
);

CREATE TABLE orders_line (
                             id BIGSERIAL NOT NULL PRIMARY KEY,
                             order_id INT REFERENCES orders (id),
                             menu_position_id INT REFERENCES menu_positions(id),
                             count INT DEFAULT 0
);

INSERT INTO orders(status) VALUES ('new');
INSERT INTO orders_line(order_id, menu_position_id, count) VALUES (1, 2, 1);
INSERT INTO orders_line(order_id, menu_position_id, count) VALUES (1, 3, 2);