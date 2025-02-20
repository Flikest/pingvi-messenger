CREATE TABLE IF NOT EXISTS users(
    id UUID UNIQUE NOT NULL,
    name VARCHAR(1000) UNIQUE NOT NULL,
    pass VARCHAR(1000) NOT NULL,
    about_me TEXT
);

CREATE TYPE IF NOT EXISTS messege (
    id BIGSERIAL NOT NULL,
    content TEXT NOT NULL,
    sending_time TIME NOT NULL
);

CREATE TYPE IF NOT EXISTS participant(
    user_id UUID NOT NULL,
    name varchar(1000) NOT NULL,
    about_me TEXT
);

CREATE TABLE IF NOT EXISTS chats(
    id UUID UNIQUE NOT NULL,
    name VARCHAR(1000) NOT NULL,
    img TEXT
    messages message,
    participants participant[] NOT NULL
);

