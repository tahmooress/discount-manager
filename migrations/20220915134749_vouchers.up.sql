BEGIN;

CREATE TABLE IF NOT EXISTS vouchers (
    id varchar(26) PRIMARY KEY NOT NULL,
    campaign_id varchar(26) NOT NULL,
    code TEXT NOT NULL,
    value decimal(64, 0) NOT NULL,
    redeemed BOOL DEFAULT FALSE,
    created_at TIMESTAMP(6) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    updated_at TIMESTAMP(6) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    CONSTRAINT vouchers_campaign_id_fk FOREIGN KEY (campaign_id) REFERENCES campaigns (id)
);

CREATE INDEX vouchers_campaign_id_redeemed_index ON vouchers (campaign_id, redeemed);

COMMIT;