CREATE TABLE users(
    "id" bigserial not null primary key,
    "login" varchar not null unique,
    "encrypted_password" varchar not null
);

CREATE TABLE LoginWithPassword(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "login" varchar,
    "password" varchar,
    "created_at" timestamp default NOW()
);

CREATE TABLE CreditCard(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "owner_name" varchar,
    "owner_last_name" varchar,
    "number" int,
    "cvc" int,
    "created_at" timestamp default NOW()
);

CREATE TABLE SecretText(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "text" text,
    "created_at" timestamp default NOW()
);

CREATE TABLE SecretBinary(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "binary" bytea,
    "created_at" timestamp default NOW()
);