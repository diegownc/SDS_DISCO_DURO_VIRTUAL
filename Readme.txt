//Instalar postgre -> https://www.postgresql.org/download/linux/
//Comandos para crear la base de datos:

su - postgres 

psql

create user root with encrypted password 'root123';
create user normaluser with encrypted password 'contra123';
create database SDS;
grant all privileges on database SDS to root;

exit

psql sds

create table users(
username character varying (255) not null,
password character varying (255) not null, 
idfolder serial,
constraint pk_user primary key (username),
constraint u_folder unique (idfolder));

grant select, insert on table users to normaluser;
grant usage, select on sequence users_idfolder_seq to normaluser;


//////////////////////////////////////////////////////7
Pasos para ejecutar este programa

-- Nos descargamos las dependencias
go mod tidy

go run Login.go

