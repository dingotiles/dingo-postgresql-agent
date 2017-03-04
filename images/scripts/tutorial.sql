drop table demo;
create table demo (value text);

insert into demo values ('row1');
insert into demo values ('row2');
select * FROM demo;

select pg_switch_xlog();
