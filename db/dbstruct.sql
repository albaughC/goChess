create table users (
    id UUID PRIMARY KEY,
    username varchar(30) NOT NULL UNIQUE,
    password text NOT NULL,
    email varchar(60) NOT NULL UNIQUE,
    createdate date NOT NULL,
    icon varchar(30),
    elo int,
    gamesplayed int
);

create table games (
    id SERIAL PRIMARY KEY,
    playedon date NOT NULL,
    wplayer UUID NOT NULL REFERENCES users,
    bplayer UUID NOT NULL REFERENCES users,
    winner UUID CHECK (winner = wplayer OR winner=bplayer),
    check (wplayer <> bplayer)
);

create table moves (
    id UUID PRIMARY KEY,
    gameid INT REFERENCES games,
    movenum int NOT NULL,
    notation varchar(10) NOT NULL,
    state JSON NOT NULL
);
