create table if not exists messages (
    uuid uuid primary key not null,
    content varchar(512) not null,
    ts  bigint not null
);