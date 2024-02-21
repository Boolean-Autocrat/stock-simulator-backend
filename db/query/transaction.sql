-- name: BeginTransaction :exec
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- name: EndTransaction :exec
COMMIT;