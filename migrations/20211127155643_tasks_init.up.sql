create table if not exists "tasks"(
  "id" bigserial primary key,
  "day_id" bigint references "days"("id"),
  "name" text not null,
  "start" time with time zone not null,
  "end" time with time zone null
);