CREATE TABLE IF NOT EXISTS users (
	id uuid primary key,
	fullname varchar(200),
	email varchar(100),
	password varchar(100),
	is_active boolean default false,
	deleted boolean default false,
	created_date timestamp
);

CREATE TABLE IF NOT EXISTS roles (
	id uuid primary key,
	name varchar(100),
	description varchar(200),
	deleted boolean default false,
	created_date timestamp
);
INSERT INTO roles(id,name,description,created_date) VALUES('3613da52-f643-4108-83db-64a114714ba8','admin','Administrator','2023-05-26 10:00:00');
INSERT INTO roles(id,name,description,created_date) VALUES('507a93b6-69b8-4ea4-b108-d165ea8f52d8','worker','Worker','2023-05-26 10:00:00');
INSERT INTO roles(id,name,description,created_date) VALUES('3bc21d7a-16ec-45a7-acca-17e52d2b954e','purchasing','Purchasing','2023-06-11 10:00:00');

CREATE TABLE IF NOT EXISTS user_roles (
	id uuid primary key,
	user_id uuid,
	role_id uuid,
	deleted boolean default false,
	created_date timestamp,
	foreign key (user_id) references users(id) on delete cascade,
	foreign key (role_id) references roles(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS branches (
	id uuid primary key,
	code varchar(100),
	name varchar(200),
	location text,
	phone varchar(15),
	email varchar(100),
	is_main boolean default false,
	is_active boolean default false,
	deleted boolean default false,
	created_date timestamp
);

CREATE TABLE IF NOT EXISTS branch_users (
	id uuid primary key,
	user_id uuid,
	branch_id uuid,
	deleted boolean default false,
	created_date timestamp,
	foreign key (user_id) references users(id) on delete cascade,
	foreign key (branch_id) references branches(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS categories (
	id uuid primary key,
	name varchar(100),
	description text,
	images jsonb default null,
	deleted boolean default false,
	created_date timestamp
);

CREATE TABLE IF NOT EXISTS products (
	id uuid primary key,
	category_id uuid,
	code varchar(30),
	name varchar(200),
	price float default 0,
	stock numeric default 0,
	with_stock boolean default false,
	description text,
	images jsonb default null,
	is_active boolean default false,
	deleted boolean default false,
	created_date timestamp,
	foreign key (category_id) references categories(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS customers (
	id uuid primary key,
	code varchar(50),
	name varchar(150),
	phone varchar(15),
	email varchar(60),
	address text,
	details jsonb default null,
	deleted boolean default false,
	created_date timestamp
);

CREATE TABLE IF NOT EXISTS trx (
	id uuid primary key,
	user_id uuid,
	customer_id uuid,
	branch_id uuid,
	code varchar(20),
	total_qty numeric default 0,
	total_price numeric default 0,
	discount numeric default 0,
	discount_desc text,
	grand_total numeric default 0,
	status char(2) default 'S1',
	note text,
	deleted boolean default false,
	created_date timestamp,
	foreign key (user_id) references users(id) on delete cascade,
	foreign key (customer_id) references customers(id) on delete cascade,
	foreign key (branch_id) references branches(id) on delete cascade
);
COMMENT ON COLUMN trx.status IS 'S1=open, S2=paid, S3=canceled';

CREATE TABLE IF NOT EXISTS trx_items (
	id uuid primary key,
	trx_id uuid,
	product_id uuid,
	qty numeric default 0,
	price numeric default 0,
	foreign key (trx_id) references trx(id) on delete cascade,
	foreign key (product_id) references products(id) on delete cascade
);

-- PURCHASE ORDER
CREATE TABLE IF NOT EXISTS suppliers (
	id uuid primary key,
	code varchar(50),
	name varchar(200),
	phone varchar(15),
	email varchar(100),
	address text,
	details jsonb default null,
	deleted boolean default false,
	created_date timestamp,
	updated_date timestamp	
);

CREATE TABLE IF NOT EXISTS purchase_orders (
	id uuid primary key,
	user_id uuid,
	branch_id uuid,
	supplier_id uuid,
	no_purchase varchar(60),
	purchase_date timestamp,
	total_qty numeric default 0,
	total_price numeric default 0,
	discount numeric default 0,
	status char(2) default 'S1',
	note text,
	deleted boolean default false,
	created_date timestamp,
	updated_date timestamp,
	foreign key (user_id) references users(id) on delete cascade,
	foreign key (branch_id) references branches(id) on delete cascade,
	foreign key (supplier_id) references suppliers(id) on delete cascade
);
COMMENT ON COLUMN purchase_orders.status IS 'S1=draft, S2=confirm, S3=paid, S4=rejected';

CREATE TABLE IF NOT EXISTS purchase_order_items (
	id uuid primary key,
	purchase_order_id uuid,
	product_id uuid,
	qty numeric default 0,
	price numeric default 0,
	foreign key (purchase_order_id) references purchase_orders(id) on delete cascade,
	foreign key (product_id) references products(id) on delete cascade
);


