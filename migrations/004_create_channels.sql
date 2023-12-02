create table channels(
  id generated always as identity primary key,
  name varchar(255),
  discord_id bigint,
  guild_id integer
);

---- create above / drop below ----

drop table channels;
