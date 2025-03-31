BEGIN;

CREATE TABLE IF NOT EXISTS public.users (
 id integer primary key generated always as identity,
 login varchar NOT NULL,
 password varchar NOT NULL,
 CONSTRAINT login UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS public.orders (
  id integer primary key,
  user_id integer NOT NULL,
  status varchar NOT NULL,
  uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL,
  accural float NOT NULL,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id)
  REFERENCES users(id)

);

CREATE TABLE IF NOT EXISTS public.withdrawals (
  id integer primary key generated always as identity,
  order_id integer,
  sum integer,
  processed_at TIMESTAMP WITH TIME ZONE NOT NULL,
  CONSTRAINT fk_order_id FOREIGN KEY (order_id)
  REFERENCES orders(id)
);


CREATE TABLE IF NOT EXISTS public.balance (
  user_id integer,
  current integer,
  withdraw integer,
  CONSTRAINT fk_user_id FOREIGN KEY (user_id)
  REFERENCES users(id)
);
COMMIT;