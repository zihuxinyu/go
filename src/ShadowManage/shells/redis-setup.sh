#!/bin/sh
yum install -y m2crypto python-setuptools gcc-c++ tcl unzip nano
easy_install pip
pip install supervisor


#redis安装
wget http://download.redis.io/redis-stable.tar.gz
tar xvzf redis-stable.tar.gz
cd redis-stable
make distclean
make
make test


#redis 配置
cd src
cp redis-server /usr/local/bin/
cp redis-cli /usr/local/bin/
mkdir /etc/redis
mkdir /var/redis
mkdir /var/redis/log
mkdir /var/redis/run
mkdir /var/redis/6379

cat >/etc/redis/redis.conf<<EOF
daemonize yes
pidfile /var/redis/run/redis_6379.pid
logfile /var/redis/log/redis_6379.log
dir /var/redis/6379
port 0
unixsocket /tmp/redis.sock
unixsocketperm 755
tcp-backlog 511
timeout 0
tcp-keepalive 0
loglevel notice
databases 16
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir ./
slave-serve-stale-data yes
slave-read-only yes
repl-disable-tcp-nodelay no
slave-priority 100
appendonly no
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
lua-time-limit 5000
slowlog-log-slower-than 10000
slowlog-max-len 128
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-entries 512
list-max-ziplist-value 64
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit slave 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
aof-rewrite-incremental-fsync yes
EOF


wget -q https://raw.github.com/saxenap/install-redis-amazon-linux-centos/master/redis-server
mv redis-server /etc/init.d
chmod 755 /etc/init.d/redis-server
chkconfig --add redis-server
chkconfig --level 345 redis-server on
####
# To start Redis just uncomment this line
####
service redis-server start