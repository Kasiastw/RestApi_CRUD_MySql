CREATE TABLE `tables` (
                           `id` INT NOT NULL auto_increment,
                           PRIMARY KEY (`id`),
                           capacity int,
                           booked_seats int,
                           available_seats int
);
CREATE TABLE `guestsList` (
                          `id` INT NOT NULL auto_increment,
                          PRIMARY KEY (`id`),
                          table_id INT,
                          name VARCHAR(100) NOT NULL,
                          accompanying_guests INT,
                          status int,
                          arrival_time bigint,
                          FOREIGN KEY (table_id) REFERENCES tables(id)
);
