-- +migrate Up

-- +migrate StatementBegin
CREATE SEQUENCE IF NOT EXISTS customer_pkey_seq;
CREATE TABLE "customer"
(
    id         BIGINT NOT NULL DEFAULT nextval('customer_pkey_seq'::regclass),
    name       VARCHAR(256),
    phone      VARCHAR(30),
    email      VARCHAR(256),
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted    BOOLEAN         DEFAULT FALSE,
    CONSTRAINT pk_customer_id PRIMARY KEY (id)
);

INSERT INTO customer
    (name, phone, email)
VALUES ('Cust1', '0858000001', 'cust1@mail.com'),
       ('Cust2', '0858000002', 'cust2@mail.com'),
       ('Cust3', '0858000003', 'cust3@mail.com'),
       ('Cust4', '0858000004', 'cust4@mail.com'),
       ('Cust5', '0858000005', 'cust5@mail.com');

CREATE SEQUENCE IF NOT EXISTS loans_pkey_seq;
CREATE TABLE "loans"
(
    id             BIGINT NOT NULL DEFAULT nextval('loans_pkey_seq'::regclass),
    customer_id    BIGINT,
    loan_amount    FLOAT8 NOT NULL DEFAULT 0,
    interest_rate  FLOAT8 NOT NULL DEFAULT 0,
    total_payable  FLOAT8 NOT NULL DEFAULT 0,
    total_weeks    INT    NOT NULL DEFAULT 0,
    weekly_payment FLOAT8 NOT NULL DEFAULT 0,
    out_standing   FLOAT8 NOT NULL DEFAULT 0,
    start_date     DATE,
    end_date       DATE,
    is_delinquent  BOOL            DEFAULT FALSE,
    created_at     TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted        BOOLEAN         DEFAULT FALSE,
    CONSTRAINT pk_loans_id PRIMARY KEY (id),
    CONSTRAINT fk_loan_customerid FOREIGN KEY (customer_id) REFERENCES customer (id)
);

CREATE SEQUENCE IF NOT EXISTS loansschedule_pkey_seq;
CREATE TABLE "loans_schedule"
(
    id          BIGINT NOT NULL DEFAULT nextval('loansschedule_pkey_seq'::regclass),
    loan_id     BIGINT,
    week_number INT    NOT NULL DEFAULT 0,
    due_date    DATE,
    amount_due  FLOAT8 NOT NULL DEFAULT 0,
    is_paid     BOOL            DEFAULT FALSE,
    created_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted     BOOLEAN         DEFAULT FALSE,
    CONSTRAINT pk_loanschedule_id PRIMARY KEY (id),
    CONSTRAINT fk_loanschedule_loanid FOREIGN KEY (loan_id) REFERENCES loans (id)
);

CREATE SEQUENCE IF NOT EXISTS payments_pkey_seq;
CREATE TABLE "payments"
(
    id           BIGINT NOT NULL DEFAULT nextval('payments_pkey_seq'::regclass),
    loan_id      BIGINT,
    payment_date DATE,
    amount_paid  FLOAT8 NOT NULL DEFAULT 0,
    created_at   TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted      BOOLEAN         DEFAULT FALSE,
    CONSTRAINT pk_payments_id PRIMARY KEY (id),
    CONSTRAINT fk_payments_loanid FOREIGN KEY (loan_id) REFERENCES loans (id)
);


-- +migrate StatementEnd
