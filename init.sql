-- alter user client createdb;
-- ALTER ROLE
CREATE USER devuser;
CREATE DATABASE storage_db;


\connect storage_db;

CREATE TABLE kv_storage (
	id	SERIAL PRIMARY KEY,
	key_name	text UNIQUE,
	key_value	text,
	datetime	timestamp DEFAULT current_timestamp
);
GRANT ALL PRIVILEGES ON DATABASE storage_db TO devuser;
GRANT ALL PRIVILEGES ON TABLE kv_storage TO devuser;
GRANT CONNECT ON DATABASE storage_db TO devuser;
GRANT USAGE ON SCHEMA public TO devuser;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO devuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO devuser;

