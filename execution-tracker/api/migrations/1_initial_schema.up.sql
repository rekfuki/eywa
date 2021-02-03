CREATE TABLE timeline_logs (
    request_id uuid NOT NULL,
    user_id uuid NOT NULL,
    function_id uuid NOT NULL,
    event_name text,
    event_type text,
    response int,
    method text default '---',
    duration bigint,
    "timestamp" timestamp without time zone,
    expires_at timestamp without time zone
);

CREATE INDEX timeline_logs_request_id_idx ON timeline_logs USING btree (request_id);
CREATE INDEX timeline_logs_user_id_idx ON timeline_logs USING btree (user_id);
CREATE INDEX timeline_logs_function_id_idx ON timeline_logs USING btree (function_id);
CREATE INDEX timeline_logs_event_type_idx ON timeline_logs USING btree (event_type);
CREATE INDEX timeline_logs_method_idx ON timeline_logs USING btree (method);
CREATE INDEX timeline_logs_response_idx ON timeline_logs USING btree (response);
CREATE INDEX timeline_logs_timestamp_idx ON timeline_logs USING btree ("timestamp");
CREATE INDEX timeline_logs_expires_at_idx ON timeline_logs USING btree (expires_at);

CREATE TABLE event_logs (
    request_id uuid NOT NULL,
    user_id uuid,
    type text,
    function_name text,
    function_id text,
    message text,
    is_error boolean default TRUE,
    "timestamp" timestamp without time zone,
    expires_at timestamp without time zone
);

CREATE INDEX event_logs_request_id_idx ON event_logs USING btree (request_id);
CREATE INDEX event_logs_user_id_idx ON event_logs USING btree (user_id);
CREATE INDEX event_logs_type_idx ON event_logs USING btree (type);
CREATE INDEX event_logs_is_error_idx ON event_logs USING btree (is_error);
CREATE INDEX event_logs_function_name_idx ON event_logs USING btree (function_name);
CREATE INDEX event_logs_function_id_idx ON event_logs USING btree (function_id);
CREATE INDEX event_logs_message_idx ON event_logs USING gin (message gin_trgm_ops);
CREATE INDEX event_logs_timestamp_idx ON event_logs USING btree ("timestamp");
CREATE INDEX event_logs_expires_at_idx ON event_logs USING btree (expires_at);