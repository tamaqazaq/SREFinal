CREATE TABLE "UserPoints" (
                              user_id INT PRIMARY KEY,
                              points INTEGER DEFAULT 0,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              FOREIGN KEY (user_id) REFERENCES "User" (user_id) ON DELETE CASCADE
);

CREATE TABLE "Achievements" (
                                id SERIAL PRIMARY KEY,
                                user_id INT,
                                achievement_name VARCHAR(255),
                                achieved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                FOREIGN KEY (user_id) REFERENCES "UserPoints" (user_id) ON DELETE CASCADE
);