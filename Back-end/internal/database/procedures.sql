DROP PROCEDURE IF EXISTS login_limit;

DELIMITER //

CREATE PROCEDURE login_limit (
    IN check_the_email VARCHAR(120),
    IN check_the_name VARCHAR(50),
    IN success BOOLEAN
)

BEGIN

    DECLARE attempts INT;

    IF success THEN
        DELETE FROM login_attempts
        WHERE (email = check_the_email OR name = check_the_name);
        SELECT FALSE AS check_the_block;

    ELSE
        INSERT INTO login_attempts (email, name, success, attempt_in)
        VALUES (check_the_email, check_the_name, 1, NOW());
        DELETE FROM login_attempts
        WHERE attempt_in < NOW() - INTERVAL 1 DAY;
        SELECT COUNT(*) INTO attempts
        FROM login_attempts
        WHERE (email = check_the_email OR name = check_the_name)
            AND success = 0
            AND attempt_in > NOW() - INTERVAL 1 HOUR

        IF attempts >= 5 THEN
            SELECT TRUE AS check_the_block;
        ELSE
            SELECT FALSE AS check_the_block;
        END IF;

END //

DELIMITER ;