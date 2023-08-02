
Create Table if not exists sessions (
                                        id bigint primary key auto_increment,
                                        created_at timestamp default current_timestamp,
                                        updated_at timestamp default current_timestamp,
                                        access_token text,
                                        user_id bigint ,
                                        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE

);
