start transaction;

--
-- 用户表
--
create table if not exists "user" (
    id serial primary key,
    name varchar(128) not null,                                     -- 用户名
    status int not null default 0,                                  -- 用户状态 0: 有效 1: 无效
    create_at timestamp not null default now(),
    update_at timestamp not null default now()
);

--
-- 认证元信息, 存储AK/SK
--
create table if not exists auth_meta (
    id serial primary key,
    user_id int not null,                                           -- 用户id 外键: user.id
    ak varchar(32) not null,
    sk varchar(32) not null,
    is_valid bool default true,                                     -- 当前ak/sk是否有效 true: 有效 false: 无效
    create_at timestamp not null default now(),
    update_at timestamp not null default now()
);

--
-- bucket(桶信息)
--
-- 说明: 一个用户有多个ak/sk, 一个用户有多个bucket.
--       如果一个用户有2个有效的ak/sk账号, 那么该用户的这两个ak/sk账号都可以访问该用户下的所有bucket.
--
create table if not exists bucket (
    id serial primary key,
    user_id int not null,                                           -- 用户id 外键: user.id
    slave_bucket_id int not null default 0,                         -- 备份桶id 外键: bucket.id 关联自身
    name varchar(64) not null,                                      -- 桶名称
    status int not null default 0,                                  -- 桶的状态: 0: 有效 1: 无效
    is_encrypted int not null default 0,                            -- 是否开启数据加密 0: 不开启 1: 开启
    acl_type int not null default 0,                                -- 访问类型 0: private 1: public-read 2: public-read-write
    obj_cache_open int not null default 1,                          -- 是否开启对象缓存 0: 不开启 1: 开启
    create_at timestamp not null default now(),
    update_at timestamp not null default now()
);

--
-- meta(存储object元信息)
--
-- 说明: 一个bucket下存储了那些object.
-- 说明: meta是分表, 分表大小为: 000 ~ 255
--
create table if not exists meta_000 (
    id serial primary key,
    user_id int not null,                                           -- 用户id 外键: user.id
    bucket_id int not null,                                         -- 桶id 外键: bucket.id
    object varchar(512) not null,                                   -- 对象名
    version varchar(32) not null,                                   -- 对象版本
    size bigint not null,                                           -- 对象大小
    user_meta varchar(512),                                         -- 用户定义的元信息
    refer_object varchar(512) not null,                             -- 引用对象名
    refer_version varchar(32) not null,                             -- 引用对象版本
    granted_r_users varchar(256) not null default '',               -- 读权限的user_id列表, 逗号分隔
    granted_rw_users varchar(256) not null default '',              -- 读写权限的user_id列表, 逗号分隔
    is_deleted bool default false,                                  -- object是否已被删除
    create_at timestamp not null default now(),
    update_at timestamp not null default now(),
    delete_at timestamp
);

-- 
-- recoup(回收)
-- 说明: 那些object做了回收
-- 
create table if not exists recoup (
    id serial primary key,
    bucket_name varchar(64),                                        -- bucket名
    object_name varchar(512),                                       -- object名
    object_version varchar(32),                                     -- object版本
    username varchar(128),                                          -- 用户名
    status int default 0,                                           -- 0: wait 1: processing 2: success 3: failed
    retry_num int default 0,                                        -- 重试次数
    create_at timestamp not null default now(),
    update_at timestamp not null default now(),
    finish_at timestamp
);


commit;
