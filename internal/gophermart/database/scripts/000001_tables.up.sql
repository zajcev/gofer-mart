BEGIN;

CREATE TABLE IF NOT EXISTS public.users (
 id integer primary key generated always as identity,
 login varchar NOT NULL,
 password varchar NOT NULL,
 CONSTRAINT login UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS public.orders (
  id varchar NOT NULL,
  user_id integer NOT NULL,
  status varchar NOT NULL,
  uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
  accural float,
  CONSTRAINT id UNIQUE (id)
);

CREATE TABLE IF NOT EXISTS public.withdrawals (
  id integer primary key generated always as identity,
  order_id varchar NOT NULL,
  sum integer,
  processed_at TIMESTAMP WITH TIME ZONE NOT NULL,
  CONSTRAINT order_id UNIQUE (order_id)
);


CREATE TABLE IF NOT EXISTS public.balance (
  user_id integer,
  current integer,
  withdraw integer
);

CREATE TABLE IF NOT EXISTS public.sessions (
  user_id integer,
  token varchar,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id)
  REFERENCES users(id),
  CONSTRAINT token UNIQUE (token)
);

COMMIT;