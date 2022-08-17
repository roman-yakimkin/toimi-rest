# Тестовый REST-сервис

Тестовый REST-сервис согласно тестовому заданию

Может осуществлять следующие действия

1. Добавить объявление 

Пример запроса
http://localhost:8080/advert POST
{
    "title" : "Advert 1",
    "description" : "Advert 1 description",
    "photos": [
    "/img/photo1",
    "/img/photo2",
    "/img/photo3"
    ],
    "price" : 12500
}

2. Изменить объявление

Пример запроса
http://localhost:8080/advert PUT
{
    "title" : "Advert 1",
    "description" : "Advert 1 description",
    "photos": [
    "/img/photo1",
    "/img/photo2",
    "/img/photo3"
    ],
    "price" : 12500
}

3. Удалить объявление

Пример запроса
http://localhost:8080/advert/1 DELETE

4. Получить объявление по ID

Пример запроса
http://localhost:8080/advert/2 GET

5. Получить страницу с объявлениями

Пример запроса
http://localhost:8080/adverts?page=0&sort=price-desc GET

Хранение данных реализовано двумя способами - в памяти и в базе данных Postgres
Под каждое хранилище написан отдельный репозиторий.

Запускать можно с помощью docker-compose up, но есть пока одна проблема, которую побороть пока не удалось.
После запуска контейнера с базой данных должен скачиваться файл, который запускает миграцию - создает пустые таблицы для postgres.
К сожалению, при запуске миграции в докере-контейнере происходит сбой.

Запускется он таким образом
migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

Как можно запустить данный проект - запустить отдельно postgres каким-либо образом, далее в файле 
config/config.yml настроить параметры доступа в базе данных, далее скачать файл миграций и запустить миграцию, 
либо выполнить данный код для postgres

CREATE TABLE IF NOT EXISTS adverts (
id serial not null unique,
title varchar(200) not null,
description varchar(1000),
created timestamp not null,
price int not null
);

CREATE INDEX ON adverts(created);
CREATE INDEX ON adverts(price);

CREATE TABLE IF NOT EXISTS adverts_photos (
id serial not null unique,
advert_id int references adverts(id) not null,
photo varchar(255) not null,
delta int not null
);

CREATE INDEX ON adverts_photos(delta);

Далее в текущей директории проект компилируется и запускается
make run

Для работы с REST-запросами можно исользовать любой REST-клиент, например, Insomnia 

