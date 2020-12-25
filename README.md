s3(Simple Storage Service)
--------------------------
jinlong

# S3设计

## AK/SK 访问策略

## 预留式缓存

    redis作为缓存, 考虑性能问题: 将<=200K的数据缓存，>200K的数据直接回源

    1、读取数据: 先读缓存, 没有: 则读取数据库, 并set缓存
    2、更新数据: 先更新数据库, delete缓存(淘汰缓存)

    提高缓存命中率:
    1) redis自身LRU策略将最近最少读取的数据失效掉
    2) 缓存过期时间expiretime设置为5-30min内的一个随机值，防止缓存同时失效，请求全部转发到后端存储，导致后端存储压力过重
    3) expiretime-now < 1min时，s3触发回源异步刷新数据到redis中
    4) 为防止缓存穿透，当object在后端存储不存在时，就以这个obejct名称为前缀设置一个标识key，如：<object_name>_not_exist;
       从redis中查询数据时，先检查该标识key是否存在，如果存在，则直接返回，不需要回源到后端存储.
       不过需要注意的是，当该obejct名称有数据写入后，需要同时从redis中将该标识key删除。

## 分表策略

    1、对bucket中的object名称做hash，按hash值后两位进行分表，即分到256张表中。
