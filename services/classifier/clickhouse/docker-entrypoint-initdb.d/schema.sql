create table log (
    dt DateTime,
    id FixedString(33),
    class_id UUID,
    msg String
) engine = MergeTree()
    order by (dt, id);

create table class (
    id UUID,
    name String,
    rules Array(String)
) engine = PostgreSQL('postgres:5432', 'tenderr', 'class', 'tenderr', 'password');