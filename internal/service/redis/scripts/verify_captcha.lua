local key = KEYS[1]
local answer = ARGV[1]

local raw = redis.call('GET', key)
if not raw then return 'not_found' end
if redis.call('PTTL', key) <= 0 then return 'invalid_state' end

local state = cjson.decode(raw)
if state.value == answer then
  if not state.solved then
    state.solved = true
    redis.call('SET', key, cjson.encode(state), 'KEEPTTL')
  end
  return 'solved'
end

if state.solved then return 'invalid_value' end
state.triesLeft = state.triesLeft - 1
if state.triesLeft <= 0 then
  redis.call('DEL', key)
  return 'no_tries'
end

redis.call('SET', key, cjson.encode(state), 'KEEPTTL')
return 'invalid_value'
