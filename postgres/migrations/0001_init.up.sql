---------------------------------------------------------------------
--                        classifer scheme
---------------------------------------------------------------------
create table class (
    id uuid primary key,
    name text not null,
    rules text[]
);

INSERT INTO public.class (id, name, rules) VALUES ('e407b1e7-f1d0-4e70-819b-bccaaa2227e1', 'Клиентские ошибки', '{(?i)регистрац,"Unable to upload files","Unable to get data","Execution Timeout Expired.  The timeout period elapsed prior to completion of the operation or the server is not responding.","Невозможно скачать файл по ссылке"}');
INSERT INTO public.class (id, name, rules) VALUES ('582dcc44-2f3f-4376-b75a-8d890710096d', 'Программные ошибки', '{(?i)exception,(?i)(SELECT|WHERE|DELETE|INSERT|UPDATE),"Unable to get integration token","Unhandled error during","Execution error, job","Unable to send CreateAuction in external system","Unable to SendResult CreateAuction"}');
INSERT INTO public.class (id, name, rules) VALUES ('d80d0d43-5ae6-494f-b94a-76e36278e5b4', 'Инфраструктурные ошибки', '{"R-FAULT rabbitmq","Файл с fileIdю.+отсутствует на диске","Execution error, job"}');
INSERT INTO public.class (id, name, rules) VALUES ('00000000-0000-0000-0000-000000000000', 'Неизвестная ошибка', NULL);

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

