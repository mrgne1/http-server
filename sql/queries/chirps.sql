-- name: CreateChirp :one 
insert into chirps (id, created_at, updated_at, body, user_id)
values ($1, now(), now(), $2, $3)
returning *;

-- name: GetAllChirps :many
select
    *
from chirps
order by created_at asc;

-- name: GetChirp :one
select
    *
from chirps
where id = $1;

-- name: DeleteChirp :one
delete from chirps
where id = $1
returning *;
