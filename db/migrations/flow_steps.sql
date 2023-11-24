create table flow_steps
(
    id      serial primary key,
    name    text,
    flow_id integer references flows,
    "order" integer,
    details text,
    options text
);
