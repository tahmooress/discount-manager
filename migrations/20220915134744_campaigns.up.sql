BEGIN;

CREATE TABLE IF NOT EXISTS campaigns (
    id varchar(26) PRIMARY KEY NOT NULL,
    name varchar(200) PRIMARY KEY NOT NULL,
    status BOOL  DEFAULT FALSE NOT NULL,
    start_date TIMESTAMP(6) WITH TIME ZONE NOT NULL,
    expire_date TIMESTAMP(6) WITH TIME ZONE NULL,
    created_at TIMESTAMP(6) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);

CREATE INDEX campaigns_name_index ON campaigns (name);
CREATE INDEX campaigns_status_index ON campaigns (active);

ALTER TABLE campaigns ADD CONSTRAINT campaigns_name_uq UNIQUE (name);

COMMIT;