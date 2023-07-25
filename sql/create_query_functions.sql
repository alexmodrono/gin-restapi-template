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

-- ======== QUERY FUNCTIONS ========
-- ===== FILTER QUERIES =====
-- This fuction returns username, email, and created_at values
-- for the given input user id
CREATE OR REPLACE FUNCTION auth.get_user_by_id(for_id int)
    RETURNS TABLE
            (
                id          integer,
                username    varchar,
                email       varchar,
                password    varchar,
                created_at  date
            )
    language plpgsql
AS
$$
BEGIN
    RETURN QUERY
        SELECT u.id, u.username, u.email, u.password, u.created_at
        FROM auth.user u
        WHERE u.id = for_id;
END
$$;

-- This fuction returns username, email, and created_at values
-- for the given input user email
CREATE OR REPLACE FUNCTION auth.get_user_by_email(for_email varchar)
    RETURNS TABLE
            (
                id          integer,
                username    varchar,
                email       varchar,
                password    varchar,
                created_at  date
            )
    language plpgsql
AS
$$
BEGIN
    RETURN QUERY
        SELECT u.id, u.username, u.email, u.password, u.created_at
        FROM auth.user u
        WHERE u.email = for_email;
END
$$;

-- This fuction returns username, email, and created_at values
-- for the given input username
CREATE OR REPLACE FUNCTION auth.get_user_by_username(for_username varchar)
    RETURNS TABLE
            (
                id          integer,
                username    varchar,
                email       varchar,
                password    varchar,
                created_at  date
            )
    language plpgsql
AS
$$
BEGIN
    RETURN QUERY
        SELECT u.id, u.username, u.email, u.password, u.created_at
        FROM auth.user u
        WHERE u.username = for_username;
END
$$;