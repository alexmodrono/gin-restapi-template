/*
File Name: create_query_functions.sql
Abstract: This file contains functions that provide a convenient way
to interact with the database as they encapsulate common operations
and help reduce the complexity and length of the queries when using
the database driver.

Author: Alejandro Modroño <alex@sureservice.es>
Created: 07/10/2023
Last Updated: 07/12/2023
*/

-- ======== SCHEMAS ========
CREATE SCHEMA auth;

-- ======== TABLES ========
CREATE TABLE auth.user
(
    -- ======== KEYS ========
    id            SERIAL        not null
            primary key,
    username      varchar(100)  not null,
    email         varchar(100)  not null,
    password      varchar(100)  not null,
    created_at    date          not null,

    -- ======== CONSTRAINTS ========
    CONSTRAINT user_email_unique UNIQUE (email),
    CONSTRAINT user_username_unique UNIQUE (username)
);

ALTER TABLE auth.user
    owner to api;
