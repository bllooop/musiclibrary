-- +goose Up
-- +goose StatementBegin
CREATE TABLE songlist
(
    id serial not null unique,
    name varchar(150) not null,
    artist varchar(150) not null,
    releasedate varchar(11) not null,
    text varchar(1000) not null,
    link varchar(100) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE songlist;
-- +goose StatementEnd
