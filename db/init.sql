CREATE DATABASE pet_store; 

USE pet_store; 

CREATE TABLE category 
  ( 
     id   BIGINT auto_increment PRIMARY KEY, 
     name VARCHAR(10) NOT NULL 
  ); 

CREATE TABLE tag 
  ( 
     id   BIGINT auto_increment PRIMARY KEY, 
     name VARCHAR(10) NOT NULL 
  ); 

CREATE TABLE pet 
  ( 
     id         BIGINT auto_increment PRIMARY KEY, 
     categoryid BIGINT, 
     name       VARCHAR(10) NOT NULL, 
     tagid      BIGINT, 
     status     ENUM ('available', 'pending', 'sold'), 
     FOREIGN KEY (tagid) REFERENCES tag(id), 
     FOREIGN KEY (categoryid) REFERENCES category(id) 
  ); 

CREATE TABLE photourl 
  ( 
     url     VARCHAR(2048) NOT NULL, 
     md5hash CHAR(32) NOT NULL, 
     petid   BIGINT NOT NULL, 
     PRIMARY KEY (md5hash, petid),
     FOREIGN KEY (petid) REFERENCES pet(id) 
  ); 

CREATE TABLE `order` 
  ( 
     id       BIGINT auto_increment PRIMARY KEY, 
     petid    BIGINT, 
     quantity INT,
     shipdate DATE, 
     status   ENUM ('placed', 'approved', 'delivered'), 
     complete BOOLEAN, 
     FOREIGN KEY (petid) REFERENCES pet(id) 
  ); 

CREATE TABLE user 
  ( 
     id         BIGINT auto_increment PRIMARY KEY, 
     username   VARCHAR(20), 
     password   VARCHAR(20),
     firstname  VARCHAR(10), 
     lastname   VARCHAR(20), 
     email      VARCHAR(20), 
     phone      VARCHAR(20), 
     userstatus INT 
  ); 