BEGIN;

CREATE SEQUENCE cart_after_mock
START 11
INCREMENT 1;

CREATE SEQUENCE item_after_mock
START 11
INCREMENT 1;

CREATE TABLE IF NOT EXISTS Cart (
    Id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS Cart_item(
    id SERIAL PRIMARY KEY,
    product VARCHAR(50) NOT NULL,
    quantity INT NOt NULL,
    cart_id INT NOT NULL,
    CONSTRAINT fk_cart FOREIGN KEY (cart_id)
        REFERENCES Cart(id) 
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- Insert mock data into the Cart table
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;
INSERT INTO Cart DEFAULT VALUES ;

-- Insert mock data into the Cart_item table
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Shoes', 10, 1);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Shoes', 10, 1);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Shoes', 10, 1);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Shirt', 5, 1);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Pants', 3, 2);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Hat', 7, 2);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Socks', 20, 3);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Jacket', 2, 4);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Gloves', 4, 5);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Belt', 6, 6);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Scarf', 8, 7);
INSERT INTO Cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, 'Backpack', 1, 8);
COMMIT;