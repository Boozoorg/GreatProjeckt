-- CREATE TABLE account (
--     id BIGSERIAL PRIMARY KEY,
--     name TEXT UNIQUE NOT NULL,
--     password TEXT NOT NULL
-- )

-- CREATE TABLE account_token(
--     token TEXT  NOT NULL UNIQUE,
--     account_id BIGINT NOT NULL REFERENCES account,
--     expire TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
--     created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- )

-- CREATE TABLE messanger(
--     sendler BIGINT REFERENCES account,
--     receiver BIGINT REFERENCES account,
--     message TEXT NOT NULL,
--     time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- )

SELECT * FROM messanger

-- SELECT * FROM messanger WHERE (sendler = 1 AND receiver = 2) OR (sendler = 2 AND receiver = 1)