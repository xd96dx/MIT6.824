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
