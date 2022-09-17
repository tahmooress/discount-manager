BEGIN;

CREATE TABLE IF NOT EXISTS redeemed (
    id varchar(26) PRIMARY KEY NOT NULL,
    voucher_id varchar(26) NOT NULL,
    mobile varchar(11) NOT NULL,
    created_at TIMESTAMP(6) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    CONSTRAINT redeemed_voucher_id_fk FOREIGN KEY (voucher_id) REFERENCES vouchers (id)
);

ALTER TABLE redeemed ADD CONSTRAINT redeemed_mobile_uq UNIQUE (mobile);

COMMIT;