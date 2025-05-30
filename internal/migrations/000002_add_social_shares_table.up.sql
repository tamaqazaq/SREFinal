
CREATE TABLE SocialShares (
                              id SERIAL PRIMARY KEY,
                              campaign_id INT NOT NULL REFERENCES "Campaign"(campaign_id),
                              user_id INT NOT NULL REFERENCES "User"(user_id),
                              shared_to VARCHAR(255) NOT NULL, -- e.g., "facebook", "twitter"
                              created_at TIMESTAMP DEFAULT NOW()
);

