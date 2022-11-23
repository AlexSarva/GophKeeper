package storagepg

const ddl = `
create table if not exists public.creds (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    title text not null,
    login text not null,
    passwd text not null,
    notes text not null,
    created timestamp default now(),
    changed timestamp
);

create table if not exists public.notes (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null,
    title text not null,
    note text not null,
    created timestamp default now(),
    changed timestamp
);

create table if not exists public.files (
  id uuid primary key default gen_random_uuid(),
  user_id uuid  not null,
  title text  not null,
  file bytea  not null,
  notes text,
  created timestamp default now(),
  changed timestamp
);

create table if not exists public.cards (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null,
  title text not null,
  card_number text not null,
  card_owner text not null,
  card_exp text not null,
  notes text,
  created timestamp default now(),
changed timestamp
);`
