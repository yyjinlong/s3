start transaction;


--
-- 创建函数, 返回meta分表
--
create or replace function meta(idx text)
returns text as $$
declare
    sql text;
begin
    sql := 'create table if not exists meta_' || idx || '(
        id serial primary key,
        user_id int not null,                                           -- 用户id 外键: user.id
        bucket_id int not null,                                         -- 桶id 外键: bucket.id
        object varchar(512) not null,                                   -- 对象名
        version varchar(32) not null,                                   -- 对象版本
        size bigint not null,                                           -- 对象大小
        user_meta varchar(512),                                         -- 用户定义的元信息
        refer_object varchar(512) not null,                             -- 引用对象名
        refer_version varchar(32) not null,                             -- 引用对象版本
        granted_r_users varchar(256) not null default '''',               -- 读权限的user_id列表, 逗号分隔
        granted_rw_users varchar(256) not null default '''',              -- 读写权限的user_id列表, 逗号分隔
        is_deleted bool default false,                                  -- object是否已被删除
        create_at timestamp not null default now(),
        update_at timestamp not null default now(),
        delete_at timestamp
    );';
    return sql;
end;
$$ language plpgsql;


--
-- 定义创建meta分表的存储过程
--
create or replace function create_meta_sub_table()
returns void as $$
declare
    idx text;
    prefix text;
    sql text;
begin
    for i in 1..255 loop
        idx := ''||i;
        prefix := lpad(idx, 3, '0');
        sql := meta(prefix);
        execute sql;
    end loop;
end;
$$ language plpgsql;

select create_meta_sub_table();


commit;
