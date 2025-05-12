create table key_val
(
    key TEXT   PRIMARY KEY,
    str TEXT   DEFAULT '',
    num BIGINT DEFAULT 0
);

CREATE TABLE bot_chats
(
    chat_id BIGINT PRIMARY KEY,
    joined  TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    title   TEXT,
    link    TEXT
);

create table tasks
(
    id BIGSERIAL PRIMARY KEY,

    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,

    "hidden"   BOOLEAN DEFAULT false,
    "weight"   BIGINT DEFAULT 0,
    max_clicks BIGINT DEFAULT 0,

    "type"  TEXT,
    "name"  TEXT,
    "desc"  TEXT,
    icon    TEXT,
    premium BOOLEAN,
    points  BIGINT,

    interval BIGINT,
    pending  BIGINT,

    action_link TEXT,
    action_chat_id BIGINT,
    action_ton_amount DECIMAL,

    action_stars_amount BIGINT,
    action_stars_title  TEXT,
    action_stars_desc   TEXT,
    action_stars_item   TEXT,

    action_partner_hook TEXT DEFAULT '',
    action_partner_match TEXT DEFAULT '',
    action_tapp_ads_token TEXT,

    action_ads_gram_block_id TEXT
);


create table stars_transactions
(
    id       BIGSERIAL PRIMARY KEY,
    ts       TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    tx_id    TEXT,
    task_id  BIGINT,
    user_id  BIGINT,
    amount   BIGINT
);

CREATE UNIQUE INDEX stars_transactions_tx_id ON stars_transactions (tx_id);


create table ton_transactions
(
    id          BIGSERIAL PRIMARY KEY,
    ts          TIMESTAMP WITHOUT TIME ZONE,
    user_id     BIGINT,
    task_id     BIGINT,
    ton_amount  BIGINT,
    is_deposit  BOOLEAN,
    comment     TEXT,
    src_addr    TEXT,
    dst_addr    TEXT,
    "hash"      TEXT
);

CREATE UNIQUE INDEX unique_ton_transactions_hash ON ton_transactions ("hash");


create table users
(
    id     BIGINT PRIMARY KEY,
    ref_id BIGINT DEFAULT 0,
    joined TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    inited TIMESTAMP WITHOUT TIME ZONE,

    points     BIGINT DEFAULT 0,
    ref_points BIGINT NOT NULL DEFAULT 0,
    ref_count  BIGINT NOT NULL DEFAULT 0,

    ip_address TEXT DEFAULT '',
    user_agent TEXT DEFAULT '',

    first_name     TEXT DEFAULT '',
    last_name      TEXT DEFAULT '',
    username       TEXT DEFAULT '',
    lang_code      TEXT DEFAULT '',
    is_premium     BOOLEAN DEFAULT false,
    photo_url      TEXT DEFAULT '',
    allow_messages BOOLEAN DEFAULT false
);

CREATE INDEX idx_users_ref_id ON users (ref_id);


CREATE TYPE task_status AS ENUM ('active', 'pending', 'claim', 'done');

create table tasks_state
(
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    task_id BIGINT NOT NULL,
    sub_id  BIGINT NOT NULL,

    finish_count BIGINT NOT NULL DEFAULT 0,
    total_points BIGINT NOT NULL DEFAULT 0,

    "status" task_status NOT NULL,
    updated  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX idx_tasks_state_task ON tasks_state (task_id);
CREATE UNIQUE INDEX unique_tasks_state ON tasks_state (user_id, task_id, sub_id);


create table tasks_finished
(
    id       BIGSERIAL PRIMARY KEY,
    ts       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    user_id  BIGINT NOT NULL,
    task_id  BIGINT NOT NULL,
    sub_id   BIGINT NOT NULL,
    points   BIGINT NOT NULL
);

CREATE INDEX idx_tasks_finished_ts ON tasks_finished (ts);


create table tasks_partner_events
(
    id       BIGSERIAL PRIMARY KEY,
    ts       TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    user_id  BIGINT,
    task_id  BIGINT,
    sub_id   BIGINT,
    payout   DECIMAL DEFAULT 0
);

CREATE UNIQUE INDEX unique_tasks_partner_events ON tasks_partner_events (user_id, task_id, sub_id, ts);


create table users_refs
(
    id      BIGSERIAL PRIMARY KEY,
    from_id BIGINT NOT NULL,
    to_id   BIGINT NOT NULL,
    points  BIGINT NOT NULL,
    "level" SMALLINT NOT NULL DEFAULT 0,
);

CREATE UNIQUE INDEX unique_users_refs ON users_refs (from_id, to_id);


create table products
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE,
    deleted_at TIMESTAMP WITHOUT TIME ZONE,
    "name"     TEXT NOT NULL,
    price      BIGINT NOT NULL,
    "type"     TEXT NOT NULL,
    amount     BIGINT NOT NULL,
    badge      TEXT NOT NULL DEFAULT ''
);

create table products_tickets
(
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL,
    product_id     BIGINT NOT NULL,
    product_type   TEXT,
    product_amount BIGINT,
    claim_price    BIGINT,
    claim_at       TIMESTAMP WITHOUT TIME ZONE,
    sent_at        TIMESTAMP WITHOUT TIME ZONE,
    "status"       TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX unique_products_tickets_claim ON products_tickets (user_id, claim_at);
