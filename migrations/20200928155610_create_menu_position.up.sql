CREATE TABLE menu_positions (
                                id BIGSERIAL NOT NULL PRIMARY KEY,
                                number INT NOT NULL UNIQUE,
                                name VARCHAR NOT NULL
);

INSERT INTO menu_positions(number, name) VALUES (1, 'Pasta alla carbonara');
INSERT INTO menu_positions(number, name) VALUES (2, 'Beef stroganoff');
INSERT INTO menu_positions(number, name) VALUES (3, 'T-bone steak');
INSERT INTO menu_positions(number, name) VALUES (4, 'Tom Yam');