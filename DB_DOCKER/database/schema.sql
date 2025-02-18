CREATE TABLE IF NOT EXISTS users (
    uid INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    nom TEXT,
    prenom TEXT,
    dob TEXT,
    gender TEXT,
    preference TEXT,
    bio TEXT,
    first_step BOOLEAN DEFAULT FALSE,
    terms_accepted BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_interests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uid INTEGER,
    interest TEXT,
    FOREIGN KEY (user_uid) REFERENCES users(uid)
);

CREATE TABLE IF NOT EXISTS user_pictures (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uid INTEGER,
    picture_path TEXT,
    is_profile BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_uid) REFERENCES users(uid)
);