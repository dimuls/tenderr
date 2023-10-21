create table log (
    time DateTime,
    class_id UUID,
    id FixedString(33),
    element_id String,
    message String
) engine = MergeTree()
    order by (time, class_id, id);

create table class (
    id UUID,
    name String,
    rules Array(String)
) engine = PostgreSQL('postgres:5432', 'tenderr', 'class', 'tenderr', 'password');