create table order_items
(
    id         integer default nextval('orders_id_seq'::regclass) not null constraint orders_pkey primary key,
    product_id integer constraint orders_product_id_fkey references products on update cascade on delete cascade,
    quantity   numeric,
    order_id   integer constraint order_items_orders_id_fk references orders on update cascade on delete cascade,
    rating     integer
);
