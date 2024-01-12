# Seata-go tcc 例子

## 如何运行？

1. 先执行以下命令，启动 seata-server和zookeeper。

   ~~~shell
   cd seata-go/tcc
   docker-compose -f docker-compose.yml up -d seata-server zookeeper
   ~~~

2. 再执行 tcc/client/cmd 和 tcc/server/cmd 目录下的main()方法。