-- +goose Up
-- +goose StatementBegin
create table if not exists courier (
    id uuid primary key,
    name text not null,
    speed int not null,
    location point not null
);

create table if not exists order (
    id uuid primary key,
    courier_id uuid references courier(id),
    location point not null,
    volume int not null,
    status text not null
);

create table if not exists storage_place (
    id uuid primary key,
    order_id uuid references order(id),
    courier_id uuid references courier(id) not null,
    volume int not null,
    name text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists storage_place;
drop table if exists order;
drop table if exists courier;
-- +goose StatementEnd
