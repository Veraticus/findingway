create table guilds(
  id generated always as identity primary key,
  name varchar(255),
  discord_id bigint,
  description text
);

---- create above / drop below ----

drop table guilds;
