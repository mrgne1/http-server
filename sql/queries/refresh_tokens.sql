-- name: CreateRefreshToken :one 
insert into refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
values ($1, now(), now(), $2, $3, null)
returning *;

-- name: GetRefreshToken :one
select 
    refresh_tokens.token,
    refresh_tokens.created_at,
    refresh_tokens.updated_at,
    refresh_tokens.expires_at,
    refresh_tokens.revoked_at,
    users.email as user_email
from refresh_tokens
join users on users.id = refresh_tokens.user_id
where token = $1;

-- name: RevokeRefreshToken :one
update refresh_tokens
    set revoked_at = now()
where token = $1
returning *;

