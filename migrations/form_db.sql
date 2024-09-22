create type division as enum ('Support', 'IT', 'Billing');

alter type division owner to postgres;

create table admin
(
    id       bigserial
        primary key,
    division division not null,
    chat_id  bigint   not null
);

alter table admin
    owner to postgres;

create table appeal
(
    id          bigserial
        primary key,
    division    division                            not null,
    subject     varchar(200)                        not null,
    text        varchar(2000)                       not null,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    answered_at timestamp,
    chat_id     bigint                              not null,
    admin_id    bigint                              not null
        constraint admin_id_fkey
            references admin,
    username    varchar(32)                         not null
);

alter table appeal
    owner to postgres;