package.path = "/etc/asterisk/?.lua;./live/etc/asterisk/?.lua;/etc/asterisk/custom/?.lua;./live/etc/asterisk/custom/?.lua;/etc/asterisk/lua/?.lua;./live/etc/asterisk/lua/?.lua;./lua/?.lua;./custom/?.lua;" .. package.path
package.cpath = "/usr/lib/lua/5.4/?.so;" .. package.cpath

log = require "log"
inspect = require "inspect"
redis = require "redis"

-- TODO: refactor configurations
require "config"

local redis_conn
local redis_configured = false

function init_redis()
    logDebug("INIT REDIS")
    local ok, err = pcall(function()
        logDebug("Connecting to Redis: " .. REDIS_HOST .. ":" .. REDIS_PORT)
        redis_conn = redis.connect({
            host = REDIS_HOST,
            port = REDIS_PORT,
            password = REDIS_PASSWORD
        })

        -- check connection
        redis_conn:ping()
        redis_configured = true
    end)

    if not ok then
        logDebug("Redis connection failed: " .. tostring(err))
        redis_configured = false
        return false
    end

    return true
end


function logDebug(v)
    local m = ""

    if channel ~= nil then
        local l = channel.CDR("linkedid"):get()
        local u = channel.CDR("uniqueid"):get()
        local i
        if l ~= u then
            i = l .. ": " .. u
        else
            i = u
        end
        m = i .. ": "
    end

    m = m .. inspect(v)

    log.debug(m)
end


function handleIncomingCall(context, extension)
    app.Answer()
    local callerId = channel.CALLERID("num"):get() or "unknown"
    logDebug("handleIncomingCall | Incoming call: " .. callerId)
    local redisKey = "incoming_call_" .. callerId
    
    redis_conn:setex(redisKey, REDIS_TTL_INCOMING_CALL, os.time())
    app.Hangup()
end

if not redis_configured then
    init_redis()
end

extensions = {
    ["from-provider"] = {
        ["_X."] = handleIncomingCall
    }
}