BEGIN;

CREATE TABLE IF NOT EXISTS public.users (
 id integer primary key generated always as identity,
 login varchar NOT NULL,
 password varchar NOT NULL,
 CONSTRAINT user_login UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS public.orders (
  id varchar NOT NULL,
  user_id integer NOT NULL,
  status varchar NOT NULL,
  uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
  accural float,
  CONSTRAINT order_id UNIQUE (id)
);

CREATE TABLE IF NOT EXISTS public.withdrawals (
  id integer primary key generated always as identity,
  user_id integer NOT NULL,
  order_id varchar NOT NULL,
  sum float,
  processed_at TIMESTAMP WITH TIME ZONE NOT NULL,
  CONSTRAINT withdrawals_order_id UNIQUE (order_id)
);


CREATE TABLE IF NOT EXISTS public.balance (
  user_id integer,
  current float,
  withdraw float DEFAULT 0,
  CONSTRAINT balance_user_id UNIQUE (user_id)
);

CREATE TABLE IF NOT EXISTS public.sessions (
  user_id integer,
  token varchar,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id)
  REFERENCES users(id),
  CONSTRAINT user_token UNIQUE (token)
);

COMMIT;