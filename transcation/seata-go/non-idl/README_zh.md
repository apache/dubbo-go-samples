# Seata-go tcc 例子

## 如何运行？

1. 先执行以下命令，启动 seata-server。

   ```shell
   cd dockercompose
   docker-compose -f docker-compose.yml up -d seata-server
   ```

2. 再执行 tcc/client/cmd 和 tcc/server/cmd 目录下的 main()方法。
