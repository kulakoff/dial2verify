log.outfile = "/var/log/asterisk/pbx_lua.log"
trunk = "first"
lang = "ru"

REDIS_HOST = os.getenv("REDIS_HOST") or "redis"
REDIS_PORT = os.getenv("REDIS_PORT") or 6379
REDIS_PASSWORD = os.getenv("REDIS_PASSWORD") or nil
REDIS_TTL_INCOMING_CALL = os.getenv("REDIS_TTL_INCOMING_CALL") or 300
