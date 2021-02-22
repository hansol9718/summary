create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(255) not null,
    passhash binary(60) not null,
    username varchar(255) not null,
    firstname varchar(64) not null,
    lastname varchar(128) not null, 
    photourl varchar(500) not null
);

CREATE UNIQUE INDEX uni_email
ON users(email);

CREATE UNIQUE INDEX uni_username
ON users(username);