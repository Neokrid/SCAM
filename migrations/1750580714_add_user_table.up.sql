DROP TABLE IF EXISTS roles;
CREATE TABLE IF NOT EXISTS roles(
    id int primary key,
    name varchar not null,
    permissions jsonb
);

INSERT INTO roles VALUES
(1, 'ADMINT'),
(2, 'CLIENT');

CREATE TABLE IF NOT EXISTS users(
    id uuid primary key,
    username varchar unique not null,
    full_name varchar not null,
    birth_date DATE,
    status varchar default '',
    email varchar unique not null,
    password varchar not null,
    img_url varchar not null DEFAULT '',
    confirmed_email boolean default 'f',
    last_seen_at TIMESTAMP DEFAULT NOW(),
    created_at timestamptz not null
);