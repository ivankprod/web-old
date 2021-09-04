#!/usr/bin/env tarantool

local log = require('log')

-- DB structure init
local function init()
	box.schema.user.create('operator', { password = '83467738438', if_not_exists = true })
	box.schema.user.grant('operator', 'read,write,execute', 'universe', nil, { if_not_exists = true })

	box.schema.sequence.create('users_id_seq', { min=1, start=1, if_not_exists = true })

	local users_roles_space = box.schema.space.create('users_roles', { if_not_exists = true })

	users_roles_space:format({
		{ name = 'id',   type = 'integer' },
		{ name = 'role', type = 'string'  },
		{ name = 'sort', type = 'integer' }
	})

	users_roles_space:create_index('primary_id', {
		if_not_exists = true,
		type = 'TREE',
		unique = true,
		parts = {{ field = 1, type = 'integer' }}
	})

	local users_types_space = box.schema.space.create('users_types', { if_not_exists = true })

	users_types_space:format({
		{ name = 'id',   type = 'integer' },
		{ name = 'type', type = 'string'  }
	})

	users_types_space:create_index('primary_id', {
		if_not_exists = true,
		type = 'TREE',
		unique = true,
		parts = {{ field = 1, type = 'integer' }}
	})

	local users_space = box.schema.space.create('users', { if_not_exists = true })

	users_space:format({
		{ name = 'user_id', type = 'integer' },
		{ name = 'user_group', type = 'integer' },
		{ name = 'user_social_id', type = 'string' },
		{ name = 'user_access_token', type = 'string' },
		{ name = 'user_avatar_path', type = 'string' },
		{ name = 'user_email', type = 'string' },
		{ name = 'user_name_first', type = 'string' },
		{ name = 'user_name_last', type = 'string' },
		{ name = 'user_last_access', type = 'string' },
		{ name = 'user_role', type = 'integer' },
		{ name = 'user_type', type = 'integer' }
	})

	users_space:create_index('primary_id', {
		sequence = 'users_id_seq',
		if_not_exists = true,
		type = 'TREE',
		unique = true,
		parts = {{ field = 1, type = 'integer' }}
	})

	users_space:create_index('secondary_group', {
		if_not_exists = true,
		type = 'TREE',
		unique = false,
		parts = {{ field = 2, type = 'integer' }}
	})

	users_space:create_index('secondary_socialid_type', {
		if_not_exists = true,
		type = 'TREE',
		unique = false,
		parts = {{ field = 3, type = 'string' }, { field = 11, type = 'integer' }}
	})
end

-- Default values
local function default_data()
	local users_roles_space = box.space.users_roles
	local users_types_space = box.space.users_types
	local users_space       = box.space.users

	users_roles_space:insert{ 1, 'Заблокирован', 4 }
	users_roles_space:insert{ 2, 'Гость', 3 }
	users_roles_space:insert{ 3, 'Вебмастер', 1 }
	users_roles_space:insert{ 4, 'Администратор', 2 }

	users_types_space:insert{ 0, 'ВКонтакте' }
	users_types_space:insert{ 1, 'Яндекс' }
	users_types_space:insert{ 2, 'Facebook' }
	users_types_space:insert{ 3, 'Google' }

	users_space:insert{
		nil, 1, 'some-social-id', 'some-access-token',
		'some-avatar-path', 'some-email', 'some-firstname', 'some-lastname',
		'some-last-access', 2, 0
	}
end

box.cfg{
	checkpoint_interval = 3600,
	checkpoint_count    = 10,

	listen     = 3301,
	pid_file   = nil,
	background = false,
	log_level  = 5
}

box.once('init-v1.4.1', init)
box.once('def-data-v1.4.1', default_data)
