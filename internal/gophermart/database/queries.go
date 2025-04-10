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
	addOrder  = "insert into orders (id, user_id, status, uploaded_at, accural) values ($1, $2, $3, $4, $5)"
	getOrder  = "select id,user_id from orders where id = $1"
	getOrders = "select id,status,uploaded_at, accural from orders where user_id = $1"
)

const (
	getBalance     = "select * from balance where user_id = $1"
	getWithdrawals = "select order_id,sum,processed_at from withdrawals where user_id = $1"
	setWithdrawals = "insert into withdrawals (order_id,user_id,sum,processed_at) values ($1, $2, $3, $4)"
)
