create table orders
(
    user_id  integer constraint orders_users_id_fk references users on update cascade on delete cascade,
    id       integer default nextval('orders_id_seq1'::regclass) not null constraint orders_pkey1 primary key,
    feedback text,
    score    integer
);
