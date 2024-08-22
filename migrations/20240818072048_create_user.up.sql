CREATE TYPE user_role AS enum ('super_admin', 'content_admin', 'user');

CREATE TABLE IF NOT EXISTS users(
    id bigserial not null primary key,
    username varchar not null unique,
    email varchar not null unique,
    role user_role not null,
    encrypted_password varchar not null
);
