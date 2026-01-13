from mysql:5.7.25
maintainer mandarava
#run mv /etc/mysql/conf.d/my.conf /etc/mysql/conf.d/my.conf.bak
run echo "Asia/Shanghai" > /etc/timezone \
&& cp -a /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 

add my.cnf /etc
user root
#mkdir -p /var/mysql/logs && mkdir -p /var/mysql/data
#docker run -p 3306:3306 --name mymysql -v /var/mysql/logs:/logs -v /var/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d registry.cn-hangzhou.aliyuncs.com/mandarava/mysql-docker:latest
#进入容器
#docker exec -it 833 /bin/bash
#允许远程访问
#use mysql;
#GRANT ALL PRIVILEGES ON *.* TO root@"%" IDENTIFIED BY "root"; 
#flush privileges;  
#select host,user from user; //查看修改是否成功。

