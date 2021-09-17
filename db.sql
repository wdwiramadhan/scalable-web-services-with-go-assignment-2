CREATE TABLE IF NOT EXISTS orders (
    order_id int NOT NULL AUTO_INCREMENT,
    customer_name varchar(100) NOT NULL,
    ordered_at timestamp NOT NULL,
    PRIMARY KEY (order_id)
);
CREATE TABLE IF NOT EXISTS items (
    item_id int NOT NULL AUTO_INCREMENT,
    item_code varchar(16) NOT NULL,
    description varchar(255)  NOT NULL,
    quantity int NOT NULL,
    order_id int NOT NULL,
    PRIMARY KEY (item_id),
    FOREIGN KEY (order_id)
        REFERENCES orders(order_id)
        ON DELETE CASCADE
);