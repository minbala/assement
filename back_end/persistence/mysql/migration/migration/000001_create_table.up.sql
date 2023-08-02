Create Table if not exists users (
                                     id bigint primary key auto_increment,
                                     created_at timestamp default current_timestamp,
                                     updated_at timestamp default current_timestamp,
                                     email varchar(255) unique not null ,
                                     name text,
                                     password text not null ,
                                     user_role varchar(20) not null
);
