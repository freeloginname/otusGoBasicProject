-- name: CreateUser :one
insert into  users (name, password)
VALUES ($1, $2)
returning id;

-- name: GetUserById :one
select * from users where id = $1;

-- name: GetUserByName :one
select * from users where name = $1;

-- name: GetAllUsers :many
select * from users;

-- name: CreateNote :one
insert into notes (user_id, name, text)
values ($1, $2, $3)
returning id;

-- name: GetNote :one
select * from notes where id = $1;

-- name: GetUserNoteByName :one
select * from notes where user_id = $1 and name = $2;

-- name: UpdateUserNoteByName :exec
update notes 
set text = sqlc.arg(text)
where name = $1 and user_id = $2
returning id;

-- name: DeleteNoteById :exec
delete from notes where id = $1
returning id;

-- name: GetUserNotes :many
select * from notes where user_id = $1;

-- name: GetAllNotes :many
select * from notes;