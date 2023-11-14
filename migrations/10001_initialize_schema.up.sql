CREATE TABLE IF NOT EXISTS users(
    "id" bigserial not null primary key,
    "login" varchar not null unique,
    "encrypted_password" varchar not null
);

CREATE TABLE IF NOT EXISTS LoginWithPassword(
    "id" bigserial not null primary key,
    "user_id" int not null,
    "name" varchar,
    "meta" text,
    "login" varchar not null,
    "password" varchar not null,
    "created_at" timestamp default NOW()
);

CREATE TABLE IF NOT EXISTS CreditCard(
    "id" bigserial not null primary key,
    "user_id" int,
    "name" varchar,
    "meta" text,
    "owner_name" varchar not null,
    "owner_last_name" varchar not null,
    "number" varchar not null,
    "cvc" varchar not null,
    "created_at" timestamp default NOW()
);

CREATE TABLE IF NOT EXISTS SecretText(
    "id" bigserial not null primary key,
    "user_id" int not null,
    "name" varchar,
    "meta" text,
    "text" text not null,
    "created_at" timestamp default NOW()
);

CREATE TABLE IF NOT EXISTS SecretFile(
    "id" bigserial not null primary key,
    "user_id" int not null,
    "name" varchar,
    "meta" text,
    "path" varchar not null,
    "created_at" timestamp default NOW()
);

CREATE UNIQUE INDEX userlogin_idx ON users (login);

CREATE UNIQUE INDEX LoginWithPasswordName_idx ON LoginWithPassword (name);
CREATE UNIQUE INDEX LoginWithPasswordUserID_idx ON LoginWithPassword (user_id);
CREATE UNIQUE INDEX CreditCardName_idx ON CreditCard (name);
CREATE UNIQUE INDEX CreditCardUserID_idx ON CreditCard (user_id);
CREATE UNIQUE INDEX SecretTextName_idx ON SecretText (name);
CREATE UNIQUE INDEX SecretTextUserID_idx ON SecretText (user_id);
CREATE UNIQUE INDEX SecretFileName_idx ON SecretFile (name);
CREATE UNIQUE INDEX SecretFileUserID_idx ON SecretFile (user_id);