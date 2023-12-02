create table duties(
  id generated always as identity primary key,
  name varchar(255)
);

create table channels_duties(
  channel_id integer,
  duty_id integer,
  CONSTRAINT fk_channel
    FOREIGN KEY(channel_id)
	  REFERENCES channels(id)
  CONSTRAINT fk_duty
    FOREIGN KEY(duty_id)
	  REFERENCES duties(id)
);

---- create above / drop below ----

drop table duties;
drop_table channels_duties;
