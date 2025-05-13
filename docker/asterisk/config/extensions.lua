package.path = "/etc/asterisk/?.lua;./live/etc/asterisk/?.lua;/etc/asterisk/custom/?.lua;./live/etc/asterisk/custom/?.lua;/etc/asterisk/lua/?.lua;./live/etc/asterisk/lua/?.lua;./lua/?.lua;./custom/?.lua;" .. package.path
package.cpath = "/usr/lib/lua/5.3/?.so;" .. package.cpath

log = require "log"
inspect = require "inspect"
http = require "socket.http"
redis = require "redis"

require "config"

local redis_conn
local redis_configured = false

local function init_redis()
    local ok, err = pcall(function()
        redis_conn = redis.connect({
            host = os.getenv("REDIS_HOST") or "redis",
            port = os.getenv("REDIS_PORT") or 6379,
            password = os.getenv("REDIS_PASSWORD") or nil
        })

        -- Тестовый запрос для проверки подключения
        redis_conn:ping()
        redis_configured = true
    end)

    if not ok then
        logDebug("Redis connection failed: " .. tostring(err))
        redis_configured = false
    end
end

--redis_conn = storage.connect({
--    host = redis_server_host,
--    port = redis_server_port
--})

if redis_server_auth and redis_server_auth ~= nil then
    redis_conn:auth(redis_server_auth)
end

function logDebug(v)
    local m = ""
    m = m .. inspect(v)
    log.debug(m)
end

function handleIncomingCall(context, extension)
    logDebug("handleIncomingCall")
    app.Answer()
    local callerId = channel.CALLERID("num"):get() or "unknown"
    logDebug("Incoming call: " .. callerId)
    local redisKey = "incoming_call_" .. callerId
    redis_conn:setex(redisKey, 300, os.time())
    app.Hangup()
end

--init_redis()
if not redis_configured then
    logDebug("INIT REDIS")
    init_redis()
end

logDebug("START")

extensions = {
    ["from-provider"] = {
        ["_X."] = handleIncomingCall
    }
}