create table if not exists "days"(
    "id" bigserial primary key,
    "telegram_id" bigint not null,
    "date" date not null,
    "start" time with time zone not null,
    "end" time with time zone null,
    UNIQUE("telegram_id", "date")
)