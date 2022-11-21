package admin

// ddl tables and queries for the first initializing of database
const ddl = `
-- drop table if exists public.proxies;
-- drop table if exists public.services;
-- drop table if exists public.acl;
-- drop table if exists public.roles;
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

create table if not exists public.roles (
    id serial primary key ,
    role_name text unique
);

insert into public.roles (role_name) values ('admin') on conflict (role_name) do nothing ;
insert into public.roles (role_name) values ('maintainer') on conflict (role_name) do nothing ;
insert into public.roles (role_name) values ('developer') on conflict (role_name) do nothing ;
insert into public.roles (role_name) values ('user') on conflict (role_name) do nothing ;

-- drop table public.acl;
create table if not exists public.acl (
    user_id uuid references public.users(id),
    role_id int references public.roles(id),
    created_by uuid not null references public.users(id),
    is_del int default 0,
    deleted_by uuid references public.users(id),
    unique (user_id, role_id)
);

-- drop table public.proxies;
create table if not exists public.proxies
(
    id            bigserial
        primary key,
    host          varchar(50) not null,
    port_http     bigint      not null,
    port_socks    bigint      not null,
    username      varchar(50) not null,
    passw         varchar(50) not null,
    proxy_type    varchar(10) not null
        constraint proxy_type_cnstr
            check ((proxy_type)::text ~ 'IPv4|IPv6'::text),
    proxy_address varchar(50),
    country       varchar(50) not null,
    rent_start    timestamp   not null,
    rent_end      timestamp   not null,
    service_info  jsonb,
    is_active     integer default 1,
    status        integer default 0,
    check_time    timestamp,
    ping          numeric
);

-- drop table public.services;
create table if not exists public.services
(
    id uuid primary key,
    description text not null,
    service_type text not null,
    host varchar(50) not null,
    port bigint not null,
    filepath text not null,
    parameters jsonb,
    parameters_type text constraint proxy_type_cnstr
        check ((parameters_type)::text ~ 'FLAG|ENV|CONF'::text),
    pid int,
    status int default 0,
    edited timestamp with time zone,
    edited_by uuid references public.users(id),
    lunched timestamp with time zone,
    lunched_by uuid references public.users(id),
    stopped timestamp with time zone,
    stopped_by uuid references public.users(id),
    deleted timestamp with time zone,
    created_by uuid not null references public.users(id),
    deleted_by uuid references public.users(id),
    created timestamp with time zone default now()
);
`
