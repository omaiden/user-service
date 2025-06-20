create table "articles" (
    id         bigserial,
    title      varchar,
    content    varchar,
    author_id  varchar not null,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    primary key (id)
);
