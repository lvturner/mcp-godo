USE testdb;
CREATE TABLE IF NOT EXISTS todos (
  id int(11) NOT NULL AUTO_INCREMENT,
  title varchar(255) NOT NULL,
  completed_at datetime DEFAULT NULL,
  due_date datetime DEFAULT NULL,
  created_date datetime NOT NULL DEFAULT current_timestamp(),
  reference_id INT DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS recurrence_patterns (
  id INT AUTO_INCREMENT PRIMARY KEY,
  todo_id VARCHAR(255) NOT NULL,
  frequency VARCHAR(50) NOT NULL,
  `interval` INT NOT NULL,
  until DATETIME,
  count INT
);