box.schema.space.create('key_value')
box.schema.space.create('users')

box.space.key_value:format({
        {name = 'key', type = 'string'},
        {name = 'value', type = 'any'}
})
box.space.users:format({
        {name = 'username', type = 'string'},
        {name = 'password', type = 'string'}
})


box.space.key_value:create_index('primary', {
        type = 'hash',
        parts = {'key'}
    })
box.space.users:create_index('primary', {
        type = 'hash',
        parts = {'username'}
    })


print('Tarantool is ready and spaces are created.')
