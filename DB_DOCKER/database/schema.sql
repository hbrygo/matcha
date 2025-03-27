/*
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

-- Table pour les chatrooms (MODIFIÉE)
CREATE TABLE IF NOT EXISTS chat_rooms (
    chatroom_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,                                    -- Optionnel, pour les noms de conversations
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table pour gérer les participants d'une chatroom
CREATE TABLE IF NOT EXISTS chat_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chatroom_id INTEGER NOT NULL,
    user_uid INTEGER NOT NULL,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chatroom_id) REFERENCES chat_rooms(chatroom_id),
    FOREIGN KEY (user_uid) REFERENCES users(uid),
    UNIQUE(chatroom_id, user_uid)                 -- Un utilisateur ne peut être qu'une fois dans une chatroom
);

-- Table pour les messages
CREATE TABLE IF NOT EXISTS messages (
    message_id INTEGER PRIMARY KEY AUTOINCREMENT,
    chatroom_id INTEGER NOT NULL,
    user_uid INTEGER NOT NULL,                    -- ID de l'expéditeur
    content TEXT NOT NULL,                        -- Contenu du message
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chatroom_id) REFERENCES chat_rooms(chatroom_id),
    FOREIGN KEY (user_uid) REFERENCES users(uid)
);

-- Index pour optimiser les requêtes fréquentes
CREATE INDEX IF NOT EXISTS idx_messages_chatroom ON messages(chatroom_id);
CREATE INDEX IF NOT EXISTS idx_participants_chatroom ON chat_participants(chatroom_id);
CREATE INDEX IF NOT EXISTS idx_participants_user ON chat_participants(user_uid);
*/

CREATE TABLE IF NOT EXISTS users (
    uid INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    nom TEXT,
    prenom TEXT,
    dob TEXT,
    gender TEXT,
    bio TEXT,
    latitude TEXT,    -- Nouvelle colonne pour coordonnées géographiques
    longitude TEXT,   -- Nouvelle colonne pour coordonnées géographiques
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

-- Nouvelle table pour gérer les préférences multiples
CREATE TABLE IF NOT EXISTS user_preferences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_uid INTEGER NOT NULL,
    preference TEXT NOT NULL,
    FOREIGN KEY (user_uid) REFERENCES users(uid),
    UNIQUE(user_uid, preference)
);

-- Table pour gérer les matchs potentiels
CREATE TABLE IF NOT EXISTS potential_matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    liker_uid INTEGER NOT NULL,
    liked_uid INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (liker_uid) REFERENCES users(uid),
    FOREIGN KEY (liked_uid) REFERENCES users(uid),
    UNIQUE(liker_uid, liked_uid)
);

-- Table pour les chatrooms
CREATE TABLE IF NOT EXISTS chat_rooms (
    chatroom_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,                                    -- Optionnel, pour les noms de conversations
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Table pour gérer les participants d'une chatroom
CREATE TABLE IF NOT EXISTS chat_participants (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chatroom_id INTEGER NOT NULL,
    user_uid INTEGER NOT NULL,
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chatroom_id) REFERENCES chat_rooms(chatroom_id),
    FOREIGN KEY (user_uid) REFERENCES users(uid),
    UNIQUE(chatroom_id, user_uid)                 -- Un utilisateur ne peut être qu'une fois dans une chatroom
);

-- Table pour les messages
CREATE TABLE IF NOT EXISTS messages (
    message_id INTEGER PRIMARY KEY AUTOINCREMENT,
    chatroom_id INTEGER NOT NULL,
    user_uid INTEGER NOT NULL,                    -- ID de l'expéditeur
    content TEXT NOT NULL,                        -- Contenu du message
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chatroom_id) REFERENCES chat_rooms(chatroom_id),
    FOREIGN KEY (user_uid) REFERENCES users(uid)
);

-- Index pour optimiser les requêtes fréquentes
CREATE INDEX IF NOT EXISTS idx_messages_chatroom ON messages(chatroom_id);
CREATE INDEX IF NOT EXISTS idx_participants_chatroom ON chat_participants(chatroom_id);
CREATE INDEX IF NOT EXISTS idx_participants_user ON chat_participants(user_uid);
CREATE INDEX IF NOT EXISTS idx_user_interests ON user_interests(user_uid);
CREATE INDEX IF NOT EXISTS idx_user_pictures ON user_pictures(user_uid);
CREATE INDEX IF NOT EXISTS idx_user_preferences ON user_preferences(user_uid);
CREATE INDEX IF NOT EXISTS idx_potential_matches_liker ON potential_matches(liker_uid);
CREATE INDEX IF NOT EXISTS idx_potential_matches_liked ON potential_matches(liked_uid);