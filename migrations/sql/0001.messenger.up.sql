CREATE TABLE IF NOT EXISTS users(
    id UUID UNIQUE NOT NULL,
    name VARCHAR(1000) NOT NULL,
    pass VARCHAR(250) NOT NULL,
    avatar TEXT,
    about_me TEXT
);

CREATE TABLE IF NOT EXISTS chats(
    id UUID UNIQUE NOT NULL,
    participants UUID[] NOT NULL,
    messages JSON
);