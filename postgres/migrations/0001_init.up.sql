---------------------------------------------------------------------
--                        classifer scheme
---------------------------------------------------------------------
create table class (
    id uuid primary key,
    name text not null,
    rules text[]
);

---------------------------------------------------------------------
--                         operator scheme
---------------------------------------------------------------------

create table user_error (
    id uuid primary key,
    url text not null,
    message text not null,
    contact jsonb
);

create table error_notification (
    id uuid primary key,
    url text,
    message text not null,
    resolved bool
);

create table error_solve_waiter (
    id uuid primary key,
    error_notification_id uuid not null,
    contact jsonb
);

