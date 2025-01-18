-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE SCHEMA IF NOT EXISTS general;

-- Function to generate ulid::uuid as increment index with sorting option
CREATE OR REPLACE FUNCTION general.new_ulid() RETURNS text
AS $$
SELECT lpad(to_hex(floor(extract(epoch FROM clock_timestamp()) * 1000)::bigint), 12, '0') || encode(gen_random_bytes(10), 'hex');
$$ LANGUAGE SQL;

-- Function converts generated ulid to uuid
CREATE OR REPLACE FUNCTION general.new_uuid() RETURNS uuid
AS $$
SELECT general.new_ulid()::uuid;
$$ LANGUAGE SQL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS general.new_uuid();
DROP FUNCTION IF EXISTS general.new_ulid();

DROP SCHEMA IF EXISTS general;
DROP EXTENSION IF EXISTS pgcrypto;
-- +goose StatementEnd
