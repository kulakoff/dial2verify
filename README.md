# dial2verify

make env and edit
```shell
cp .example.env .env
```

copy libs
```shell
cd  docker/asterisk/config/lua
wget https://raw.githubusercontent.com/rxi/log.lua/refs/heads/master/log.lua
wget https://raw.githubusercontent.com/nrk/redis-lua/refs/heads/version-2.0/src/redis.lua
```

start
```shell
docker compose --profile api up -d
```