-- +goose Up
-- +goose StatementBegin
create table if not exists outbox (
    id uuid primary key,
    event_name text not null,
    event_payload jsonb not null,
    occurred_at timestamp not null,
    processed_at timestamp
);

create index if not exists idx_outbox_processed_at on outbox(processed_at);
create index if not exists idx_outbox_not_processed on outbox(processed_at) where processed_at is null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists outbox;
-- +goose StatementEnd
