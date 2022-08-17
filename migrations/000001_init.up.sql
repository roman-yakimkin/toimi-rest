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