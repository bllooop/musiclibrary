CREATE TABLE userlist
(
    id serial not null unique,
    username varchar(255) not null unique,
    password varchar(255) not null
);

CREATE TABLE songlist
(
    id serial not null unique,
    name varchar(150) not null,
    artist varchar(150) not null,
    releasedate varchar(11) not null,
    text varchar(1000) not null,
    link varchar(100) not null
);