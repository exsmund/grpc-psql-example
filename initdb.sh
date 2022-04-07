#!/bin/bash
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$POSTGRES_DB"  <<-EOSQL
     create table users (
        id serial primary key,
        email text not null,
        name text not null
     );
EOSQL