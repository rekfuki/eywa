CREATE TABLE users (
    id uuid NOT NULL,
    name text NOT NULL,
    avatar_url text,
    oauth_provider text NOT NULL,
    oauth_provider_id text NOT NULL,
    oauth_provider_email text NOT NULL,
    oauth_provider_login text NOT NULL,
    created_at timestamp without time zone,
    last_seen_at timestamp without time zone
);

CREATE INDEX users_id_idx ON users USING btree (id);
CREATE INDEX users_oauth_provider_idx ON users USING btree (oauth_provider);
CREATE INDEX users_oauth_provider_id_idx ON users USING btree (oauth_provider_id);
CREATE INDEX users_oauth_provider_email ON users USING btree (oauth_provider_email);
CREATE INDEX users_oauth_provider_login ON users USING btree (oauth_provider_login);