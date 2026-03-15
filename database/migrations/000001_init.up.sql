DROP TABLE IF EXISTS user_friends;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       gender VARCHAR(20) NOT NULL,
                       birth_date DATE NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_friends (
                              user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                              friend_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                              PRIMARY KEY (user_id, friend_id),
                              CONSTRAINT no_self_friend CHECK (user_id <> friend_id)
);

INSERT INTO users (name,email,gender,birth_date) VALUES
                                                     ('Lionel Messi','messi@example.com','male','1987-06-24'),
                                                     ('Cristiano Ronaldo','ronaldo@example.com','male','1985-02-05'),
                                                     ('Neymar Jr','neymar@example.com','male','1992-02-05'),
                                                     ('Kylian Mbappe','mbappe@example.com','male','1998-12-20'),
                                                     ('Kevin De Bruyne','debruyne@example.com','male','1991-06-28'),
                                                     ('Mohamed Salah','salah@example.com','male','1992-06-15'),
                                                     ('Robert Lewandowski','lewa@example.com','male','1988-08-21'),
                                                     ('Erling Haaland','haaland@example.com','male','2000-07-21'),
                                                     ('Luka Modric','modric@example.com','male','1985-09-09'),
                                                     ('Karim Benzema','benzema@example.com','male','1987-12-19'),
                                                     ('Harry Kane','kane@example.com','male','1993-07-28'),
                                                     ('Antoine Griezmann','griezmann@example.com','male','1991-03-21'),
                                                     ('Vinicius Junior','vinicius@example.com','male','2000-07-12'),
                                                     ('Jude Bellingham','bellingham@example.com','male','2003-06-29'),
                                                     ('Pedri Gonzalez','pedri@example.com','male','2002-11-25'),
                                                     ('Gavi','gavi@example.com','male','2004-08-05'),
                                                     ('Phil Foden','foden@example.com','male','2000-05-28'),
                                                     ('Bukayo Saka','saka@example.com','male','2001-09-05'),
                                                     ('Marcus Rashford','rashford@example.com','male','1997-10-31'),
                                                     ('Jamal Musiala','musiala@example.com','male','2003-02-26');

INSERT INTO user_friends VALUES
                             (1,3),(1,4),(1,5),(1,6),(1,7),
                             (2,3),(2,4),(2,5),(2,8),(2,9);