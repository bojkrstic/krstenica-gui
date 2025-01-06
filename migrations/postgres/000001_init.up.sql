BEGIN;

CREATE TABLE IF NOT EXISTS priests (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    city VARCHAR(100),
    title VARCHAR(100),
    status VARCHAR(100), 
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tamples (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    status VARCHAR(255),  
    city VARCHAR(255),
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS eparhije (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    status VARCHAR(255),  
    city VARCHAR(255),
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS persons (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    brief_name VARCHAR(100),
    occupation VARCHAR(255),
    religion VARCHAR(255),
    status VARCHAR(100),  
    city VARCHAR(100),
    address VARCHAR(255),
    country VARCHAR(100),
    role VARCHAR(50) CHECK (role IN ('mother', 'father', 'godfather','paroh')),
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS krstenice (
    id SERIAL PRIMARY KEY,
    book INTEGER NOT NULL,
    page INTEGER NOT NULL,
    current_number INTEGER NOT NULL,
    eparhija_id INTEGER NOT NULL REFERENCES eparhije(id),
    tample_id INTEGER NOT NULL REFERENCES tamples(id),
    parent_id INTEGER NOT NULL REFERENCES persons(id),
    godfather_id INTEGER NOT NULL REFERENCES persons(id),
    paroh_id INTEGER NOT NULL REFERENCES persons(id),
    priest_id INTEGER NOT NULL REFERENCES priests(id),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    gender VARCHAR(100) NOT NULL,
    city VARCHAR(100),
    country VARCHAR(100),
    birth_date TIMESTAMP,
    birth_order INTEGER NOT NULL,
    place_of_birthday VARCHAR(100),
    municipality_of_birthday VARCHAR(100),
    baptism TIMESTAMP,
    is_church_married BOOLEAN NOT NULL DEFAULT TRUE,
    is_twin BOOLEAN NOT NULL DEFAULT FALSE,
    has_physical_disability BOOLEAN NOT NULL DEFAULT FALSE,
    anagrafa VARCHAR(255),
    number_of_certificate INTEGER,
    town_of_certificate VARCHAR(100),
    certificate DATE,
    comment VARCHAR(255),
    status VARCHAR(255),  
    created_at TIMESTAMP
);

COMMIT;