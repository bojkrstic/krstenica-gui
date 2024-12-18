BEGIN;

CREATE TABLE IF NOT EXISTS priests (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    status VARCHAR(255), 
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tamples (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    status VARCHAR(255),  
    city VARCHAR(255),
    created_at TIMESTAMP
);


COMMIT;