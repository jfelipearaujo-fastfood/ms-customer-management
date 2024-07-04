DROP TABLE IF EXISTS customers;

CREATE TABLE IF NOT EXISTS customers (
    id varchar(255),
    document_id varchar(255),
    document_type int,
    is_anonymous boolean,
    password varchar(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);

INSERT INTO customers (id, document_id, document_type, is_anonymous, password, created_at, updated_at) 
VALUES ('19b5408e-8ee2-47d4-953b-196d41f1e367', '33344455566', 1, false, '12345678', NOW(), NOW());