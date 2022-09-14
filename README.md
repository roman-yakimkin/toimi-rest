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

Запускать можно с помощью docker-compose up

Для работы с REST-запросами можно исользовать любой REST-клиент, например, Insomnia 