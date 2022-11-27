package admin

// ddl tables and queries for the first initializing of database
const ddl = `
-- drop table if exists public.users;

create table if not exists public.users
(
    id  uuid not null primary key,
    username      text,
    email         text unique,
    passwd        text,
    is_admin bool default false,
    token         text,
    token_expires timestamp,
    created       timestamp with time zone default now()
);
`
