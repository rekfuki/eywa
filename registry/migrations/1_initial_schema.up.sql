CREATE TABLE builds (
    image_id uuid NOT NULL,
    user_id uuid NOT NULL,
    logs text,
    state text default 'building',
    created_at timestamp without time zone
);

CREATE INDEX builds_image_id_idx ON builds USING btree (image_id);
CREATE INDEX builds_user_id_idx ON builds USING btree (user_id);
CREATE INDEX builds_state_id_idx ON builds USING btree (state);
CREATE INDEX builds_timestamp_idx ON builds USING btree (created_at);

CREATE TABLE images (
    id uuid primary key,
    user_id uuid,
    registry text,
    language text,
    name text,
    version text,
    source text,
    state text,
    size int,
    created_at timestamp without time zone
);

CREATE INDEX images_id_idx ON images USING btree (id);
CREATE INDEX images_user_id_idx ON images USING btree (user_id);
CREATE INDEX images_language_idx ON images USING btree (language);
CREATE INDEX images_name_idx ON images USING btree (name);
CREATE INDEX images_version_idx ON images USING btree (version);
CREATE INDEX images_state_idx ON images USING btree (state);
CREATE INDEX images_timestamp_idx ON images USING btree ("created_at");