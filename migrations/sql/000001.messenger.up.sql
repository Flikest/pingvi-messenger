CREATE TABLE IF NOT EXISTS users(
    id UUID UNIQUE NOT NULL,
    name VARCHAR(1000) UNIQUE NOT NULL,
    pass VARCHAR(1000) NOT NULL,
    email VARCHAR(250) NOT NULL,
    avatar TEXT,
    about_me TEXT
);

CREATE TABLE IF NOT EXISTS messeges(
    chat_id UUID NOT NULL,
    message_id INT NOT NULL
    sender_id UUID,
    content BLOB NOT NULL,
    sending_time TIME NOT NULL
);

CREATE TABLE IF NOT EXISTS contacts(
    user_id UUID NOT NULL,
    contact_id UUID NOT NULL
);

CREATE TABLE IF NOT EXISTS chats(
    id UUID UNIQUE NOT NULL,
    name VARCHAR(1000) NOT NULL,
    avatar TEXT,
    unique_link_to_join TEXT
);

CREATE TABLE IF NOT EXISTS participants(
    chat_id UUID NOT NULL,
    user_id UUID NOT NULL,
    is_admin FLOAT
);