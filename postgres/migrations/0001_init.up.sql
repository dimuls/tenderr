---------------------------------------------------------------------
--                       classifier scheme
---------------------------------------------------------------------
create table class (
    id uuid primary key,
    name text not null,
    rules text[]
);

-- TODO: remove, for easy demo purposes only
INSERT INTO public.class (id, name, rules) VALUES ('e407b1e7-f1d0-4e70-819b-bccaaa2227e1', 'Клиентские ошибки', '{(?i)регистрац,"Unable to upload files","Unable to get data","Execution Timeout Expired.  The timeout period elapsed prior to completion of the operation or the server is not responding.","Невозможно скачать файл по ссылке"}');
INSERT INTO public.class (id, name, rules) VALUES ('582dcc44-2f3f-4376-b75a-8d890710096d', 'Программные ошибки', '{(?i)exception,(?i)(SELECT|WHERE|DELETE|INSERT|UPDATE),"Unable to get integration token","Unhandled error during","Execution error, job","Unable to send CreateAuction in external system","Unable to SendResult CreateAuction"}');
INSERT INTO public.class (id, name, rules) VALUES ('d80d0d43-5ae6-494f-b94a-76e36278e5b4', 'Инфраструктурные ошибки', '{"R-FAULT rabbitmq","Файл с fileIdю.+отсутствует на диске","Execution error, job"}');
INSERT INTO public.class (id, name, rules) VALUES ('00000000-0000-0000-0000-000000000000', 'Неизвестная ошибка', NULL);

---------------------------------------------------------------------
--                         operator scheme
---------------------------------------------------------------------

create table element (
    id text primary key,
    name text not null
);

insert into element (id, name) values
    ('1946d729e2ddc19eeb747ad19561f8f9', 'Форма регистрации поставщика'),
    ('7d28f96ef7eb072cba239e904e3685dd', 'Форма регистрации заказчика'),
    ('1ebe19879a13feb517f72932b6809be8', 'Форма регистрации товарной кооперации'),
    ('1505fd9adc0767a4f299407f652da87c', 'Форма оформление заказа'),
    ('a05796eb73291632430c8a026c39dc30', 'Форма создание контракта'),
    ('f915d54fed558398af87a68a66d0c847', 'Форма заказа работы'),
    ('551f5041ec2ae18deb8ce42752e161d6', 'Форма заказа комплектующих'),
    ('7dfa99641f22ec5291ca0e12c08816ee', 'Форма заказа услуги'),
    ('8c0b615fe1bf453f1739976840724cfd', 'Форма создания предложения'),
    ('a074be56b9603811c7bfe0b9caf31d2c', 'Форма поиска предложений');

create table user_error (
    id uuid primary key,
    element_id text not null,
    message text not null,
    created_at timestamp with time zone not null,
    contact jsonb
);

create table error_notification (
    id uuid primary key,
    element_id text not null,
    message text not null,
    created_at timestamp with time zone not null,
    resolved bool not null,
    resolve_message text,
    resolved_at timestamp with time zone
);

create table error_resolve_waiter (
    id uuid primary key,
    error_notification_id uuid not null,
    contact jsonb
);

