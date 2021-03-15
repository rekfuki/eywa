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

CREATE TABLE access_tokens (
    id uuid NOT NULL,
    user_id uuid not NULL,
    name text not NULL,
    token text NOT NULL,
    created_at bigint,
    expires_at bigint
);

CREATE INDEX access_token_id_idx ON access_tokens USING btree (id);
CREATE INDEX access_token_name_idx ON access_tokens USING btree (name);
CREATE INDEX access_token_user_id_idx ON access_tokens USING btree (user_id);
CREATE INDEX access_token_token_idx ON access_tokens USING btree (token);

CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE 
        data json;
        notification json;
    
    BEGIN
    
        IF (TG_OP = 'INSERT') THEN
            data = row_to_json(NEW);
        ELSEIF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE 
            RETURN NULL; 
        END IF;
        
        notification = json_build_object(
                          'action', TG_OP,
                          'data', data);
                        
        PERFORM pg_notify('access_tokens', notification::text);
        
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS access_tokens_event on access_tokens;
CREATE TRIGGER access_tokens_event
AFTER INSERT OR UPDATE OR DELETE ON access_tokens
    FOR EACH ROW EXECUTE PROCEDURE notify_event();