CREATE TABLE IF NOT EXISTS payment_event (
    user_id INT NOT NULL,
    payment_id INT AUTO_INCREMENT PRIMARY KEY,
    deposit_amount INT NOT NULL
);

CREATE TABLE IF NOT EXISTS skipped_message (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT ,
    payment_id INT ,
    deposit_amount INT 
);