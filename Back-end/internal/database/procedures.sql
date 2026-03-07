DROP PROCEDURE IF EXISTS login_limit;

DELIMITER //

CREATE PROCEDURE login_limit (
    IN check_the_email VARCHAR(120),
    OUT check_if_is_blocked BOOLEAN
)

BEGIN
    DECLARE attempts INT;
    DECLARE check_the_limit INT DEFAULT 5;
    DECLARE check_the_window INT DEFAULT 10;

    SELECT COUNT(*) INTO attempts
    FROM login_attempts
    WHERE email = check_the_email
        AND success = 0
        AND attempt_in > NOW() - INTERVAL check_the_window MINUTE;

    IF attempts >= check_the_limit THEN
        SET check_if_is_blocked = TRUE;
    ELSE
        SET check_if_is_blocked = FALSE;
        INSERT INTO login_attempts(email, success) VALUES (check_the_email, 0);
    END IF;

    DELETE FROM login_attempts WHERE attempt_in < NOW() - INTERVAL 24 HOUR;

END //

DELIMITER ;