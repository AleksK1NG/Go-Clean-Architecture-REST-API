DROP TABLE IF EXISTS users CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;
-- CREATE EXTENSION IF NOT EXISTS postgis;
-- CREATE EXTENSION IF NOT EXISTS postgis_topology;

DROP TYPE IF EXISTS SEX_T;
CREATE TYPE SEX_T AS ENUM ('male', 'female', 'other');

CREATE TABLE users
(
    user_id      UUID PRIMARY KEY                     DEFAULT uuid_generate_v4(),
    first_name   VARCHAR(32)                 NOT NULL check ( first_name <> '' ),
    last_name    VARCHAR(32)                 NOT NULL check ( last_name <> '' ),
    email        VARCHAR(64)                 NOT NULL check ( email <> '' ),
    password     VARCHAR(250)                NOT NULL CHECK ( octet_length(password) <> 0 ),
    role         VARCHAR(10)                 NOT NULL DEFAULT 'user',
    about        VARCHAR(1024)                         DEFAULT '',
    avatar       VARCHAR(512),
    phone_number VARCHAR(20),
    address      VARCHAR(250),
    city         VARCHAR(30),
    country      VARCHAR(30),
    gender       SEX_T,
    postcode     SMALLINT,
    birthday     DATE                                 DEFAULT NULL,
    created_at   TIMESTAMP                   NOT NULL DEFAULT now(),
    updated_at   TIMESTAMP                            DEFAULT current_timestamp,
    login_date   TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT current_timestamp
);