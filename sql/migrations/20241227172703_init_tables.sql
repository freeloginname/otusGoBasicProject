-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users(
 id uuid primary key default general.new_uuid(),
 name VARCHAR(255) NOT NULL,
 password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Notes(
 id uuid primary key default general.new_uuid(),
 user_id uuid references Users(id) not null,
 name VARCHAR(255) NOT NULL,
 text text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists Notes;
DROP TABLE if exists Users;
-- +goose StatementEnd
