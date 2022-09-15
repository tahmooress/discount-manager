BEGIN;

CREATE TABLE IF NOT EXISTS redeemed (
    id varchar(26) PRIMARY KEY NOT NULL,
    voucher_id varchar(26) NOT NULL,
    user varchar(11) NOT NULL,
    created_at TIMESTAMP(6) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    CONSTRAINT balances_currencies_id_fk FOREIGN KEY (currencies_id) REFERENCES currencies (id)
);

CREATE INDEX CONCURRENTLY redeemed_user_index ON redeemed (user);

COMMIT;