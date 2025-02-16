CREATE TYPE IF NOT EXISTS msg(
    user_id UUID NOT NULL,
    content JSON,
    dispatch_time TIME(),
    dispatch_date DATE,
    reactions []smiles
);
CREATE TYPE IF NOT EXISTS post(
    content TEXT,
    dispatch_time TIME(),
    dispatch_date DATE,
    reactions []smiles
);
CREATE TYPE IF NOT EXISTS sticker(
    id UUID UNIQUE NOT NULL,
    url_img JSON NOT NULL,
);

CREATE TYPE IF NOT EXISTS smiles(
    url_img JSON NOT NULL,
);
CREATE TABLE IF NOT EXISTS users(
    id UUID UNIQUE NOT NULL,
    name VARCHAR(1000) NOT NULL,
    pass VARCHAR(250) NOT NULL,
    avatar TEXT,
    about_me TEXT
);

CREATE TABLE IF NOT EXISTS chats(
    id UUID UNIQUE NOT NULL,
    name UUID UNIQUE NOT NULL,
    img JSON
    messages msg,
);

CREATE TABLE IF NOT EXISTS groops(
    id UUID UNIQUE NOT NULL,
    name VARCHAR(1000),
    messages msg,
    participants UUID[] NOT NULL,
    is_pablic BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS channels(
    id UUID UNIQUE NOT NULL,
    posts post,

);

CREATE TABLE IF NOT EXISTS stickerpack(
    id UUID UNIQUE NOT NULL,
    content sticker[]
);

