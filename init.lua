-- Ensure 'key_value' space exists

box.cfg {
    listen = 3301,
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
    print('Space "key_value" created')
else
    print('Space "key_value" already exists')
end

-- Ensure 'users' space exists
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
    print('Space "users" created')

    -- Insert default admin user if not exists
    if not box.space.users:get{'admin'} then
        box.space.users:insert{'admin', 'presale'}
        print('Default user "admin" created')
    else
        print('Default user "admin" already exists')
    end
else
    print('Space "users" already exists')
end

print('Tarantool is ready and spaces are created.')
