# GFS 论文相关知识要点

```
本文只列举个人认为重要的几处, 详细部分需阅读gfs 论文
```

## Architecture

- master 单节点, 用来做访问控制, 
- client 读写的入口
- chunkServer 存储, 每个chunk64MB, 根据master的消息, 来进行chunk操作

```
	每个chunk 有多个副本, 存在不同的chunkServer 上, 默认3个副本, chunk 的位置信息在master 的元数据上; 
	因为元数据足够小, gfs将其存在内存中, 而元数据的相关操作(创建/删除/修改) 记录在操作日志里, 当master 故障时, 备master 可以通过操作日志将操作同步, master的所有变更操作会先写入log, 再修改内存中的metadata; 
	gfs是一个弱一致性的模型, 保证的是顺序一致性, 而客户端在读取文件时, 因为弱一致性可能各个客户端读取到的不一致, 但chunkServer 会定期上报给master, 保证事件的强一致性, 最终各个client 的视图将会一致
```

## system interactions

- leases and mutation order (租约, 变更顺序)

```
	leases 是gfs 保证更新顺序一致的主要手段, 对同一个文件的写请求, master 会按顺序授予client lease, 来保证mutation order, 对于已经写入的, 客户端只能获取read lease, 当写入完成后, client 释放lease, 其他客户端获取写入权限, lease 本身有超时, 避免客户端无限占用, 客户端也需要定期renew lease
```

- data flow

```
	单个chunk 的写入, 在副本之间是线性的, 写入时会选择网络资源更空闲的先写入, 写入时, 是client 与 chunkServer 之间直接通信
```

- Atomic Record Appends

```
	追加写是以生产消费者模型完成的, 写进的数据放在队列里, 然后写入, 每次写入会返回给客户端一个新的偏移量
	各个副本不保证字节的一致, 仅保证在相同偏移量处写入相同的record
```

- snapshot

```
	使用 copy on write 的方式实现, 复制元数据, 多个快照共享相同的数据块, 不复制数据, 当有快照请求时, 先撤销其 lease, 后进入的写入需要与master 通信才能重新获取lease, 当发现该文件/目录是snapshot文件, master 将在本地创建快照文件的新块
```

## master operation

- Namespace Management and Locking

```
	锁仅使用在master 上, 可以理解为一个RWLock, 多读一写, 对于目录树也是同样的
```

- Replica Placement

```
这里主要是网络规划和机房规划, 不细说
```

- Creation, Re-replication, Rebalancing

```
	优先将副本放在利用率低的chunkServer上, 优先复制使用率较高的数据, 复制副本的过程中, 会限制资源防止副本占用挤压正常业务, master 会定期扫描来平衡副本
```

- Garbage Collection

```
	当删除文件时, 不直接删除物理文件, 而通过日志记录, 并且重命名来隐藏文件, 当 master 扫描文件系统时, 再删除这些隐藏超过三天的文件, 此时内存中相应的元数据也会被删除, 孤立的chunk 也会在chunkServer 的扫描中被删除
```

- Stale Replica Detection

```
	副本过期检测, 当master 获取一个chunk 的replicas 时, 会更新chunk 的版本号,如果某个副本在此时没有被更新版本号成功, 则被视为过期, 会在master 的定期垃圾回收中删除旧副本
```

## FAULT TOLERANCE AND DIAGNOSIS

- High Availability

```
	快速恢复, chunk 副本, master的replica
```

- Data Integrity

```
	每个chunk 都有一个32位校验和, 来保证数据的完整, 当数据与校验和不匹配时, chunkServer 将返回master 错误, master 这时会向别的chunkServer请求数据, 当新的请求成功后, 会删除错误的副本
```