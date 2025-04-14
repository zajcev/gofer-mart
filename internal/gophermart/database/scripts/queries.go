package scripts

const (
	AddUser          = "insert into users (login, password) values ($1, $2)"
	AddSession       = "insert into sessions (user_id, token) values ($1, $2) ON CONFLICT(token) DO UPDATE SET token=EXCLUDED.token"
	GetLogin         = "select login from users where login = $1"
	GetPassword      = "select password from users where login = $1"
	GetUserIDByToken = "select user_id from sessions where token = $1"
	GetUserIDByLogin = "select id from users where login = $1"
)

const (
	AddOrder           = "insert into orders (id, user_id, status, uploaded_at, accural) values ($1, $2, $3, $4, $5)"
	GetOrder           = "select id,user_id from orders where id = $1"
	GetActiveOrders    = "select id,user_id, status, accural from orders where status in ('NEW','PROCESSING')"
	GetOrders          = "select id,status,uploaded_at, accural from orders where user_id = $1"
	UpdateOrderStatus  = "update orders set status = $1 where id = $2"
	UpdateOrderAccural = "update orders set accural = $2 where id = $1"
)

const (
	GetBalance     = "select current,withdraw from balance where user_id = $1"
	SetBalance     = "INSERT INTO balance (user_id, current) VALUES ($1, $2) ON CONFLICT(user_id) DO UPDATE SET current = balance.current + EXCLUDED.current;"
	SetWithdraw    = "UPDATE balance SET current = current - $1, withdraw = withdraw + $1 WHERE user_id = $2 AND current >= $1"
	GetWithdrawals = "select order_id,sum,processed_at from withdrawals where user_id = $1"
	SetWithdrawals = "insert into withdrawals (order_id,user_id,sum,processed_at) values ($1, $2, $3, $4)"
)
