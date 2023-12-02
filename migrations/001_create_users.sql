create table users(
  id generated always as identity primary key,
  discord_id bigint,
  description text,
  notifications_enabled boolean
);

---- create above / drop below ----

drop table users;
