


Create Table if not exists user_logs (
                                         id bigint primary key auto_increment,
                                         created_at timestamp default current_timestamp,
                                         updated_at timestamp default current_timestamp,
                                         user_id bigint,
                                         method varchar(20) not null ,
                                         request_url varchar(255) not null ,
                                         service_type varchar(255) not null ,
                                         status varchar(20) not null ,
                                         error_message text,
                                         FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);


