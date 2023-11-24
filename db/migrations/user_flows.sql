create table user_flows
(
    id            serial primary key,
    user_id       integer references users on update cascade on delete cascade,
    step_id       integer constraint user_flows_flow_steps_id_fk references flow_steps,
    details       text,
    complete      boolean default false not null,
    product_id    integer constraint user_flows_products_id_fk references products,
    order_id      integer constraint user_flows_orders_id_fk references orders,
    order_item_id integer constraint user_flows_order_items_id_fk references order_items
);
