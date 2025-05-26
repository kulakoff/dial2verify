# dial2verify
A phone verification system using Asterisk PBX, Redis, and a Go-based API. It processes incoming calls, logs caller IDs in Redis, and provides an API to check phone numbers.

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
