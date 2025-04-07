package database

const (
	addUser          = "insert into users (login, password) values ($1, $2)"
	addSession       = "insert into sessions (user_id, token) values ($1, $2) ON CONFLICT(token) DO UPDATE SET token=EXCLUDED.token"
	getLogin         = "select login from users where login = $1"
	getPassword      = "select password from users where login = $1"
	getUserIdByToken = "select user_id from sessions where token = $1"
	getUserIdByLogin = "select id from users where login = $1"
)

const (
	addOrder     = "insert into orders (id, user_id, status, uploaded_at) values ($1, $2, $3, $4)"
	getOrder     = "select id from orders where user_id = $1 and where id = $2"
	getOrderUser = "select user_id from orders where id = $1"
	getOrders    = "select * from orders where user_id = $1"
)

const (
	balance = "orders"
)
