box.cfg{
    listen = 3301,
    log_level = 5,
    memtx_memory = 128 * 1024 * 1024
}

if not box.space.key_value then
    local kv_space = box.schema.space.create('key_value')

    kv_space:format({
        {name = 'key', type = 'string'},
        {name = 'value', type = 'any'}
    })
    kv_space:create_index('primary', {
        type = 'hash',
        parts = {'key'}
    })
end

if not box.space.users then
    local user_space = box.schema.space.create('users')

    user_space:format({
        {name = 'username', type = 'string'},
        {name = 'password', type = 'string'}
    })
    user_space:create_index('primary', {
        type = 'hash',
        parts = {'username'}
    })
end

-- From TS
box.space.users:put{'admin', 'presale'}

print('Tarantool is ready and spaces are created.')
