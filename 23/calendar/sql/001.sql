create table events (
    id serial primary key,
    name varchar(256),
    start_time timestamp not null,
    end_time timestamp not null
);
create index start_idx on events using btree (start_time, end_time);