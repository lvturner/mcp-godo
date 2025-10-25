USE testdb;

CREATE TABLE IF NOT EXISTS projects (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS todos (
  id int(11) NOT NULL AUTO_INCREMENT,
  title varchar(255) NOT NULL,
  completed_at datetime DEFAULT NULL,
  due_date datetime DEFAULT NULL,
  created_date datetime NOT NULL DEFAULT current_timestamp(),
  reference_id INT DEFAULT NULL,
  project_id INT DEFAULT NULL,
  PRIMARY KEY (id),
  INDEX idx_todos_project_id (project_id)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;</search>
</search_and_replace>

CREATE TABLE IF NOT EXISTS recurrence_patterns (
  id INT AUTO_INCREMENT PRIMARY KEY,
  todo_id VARCHAR(255) NOT NULL,
  frequency VARCHAR(50) NOT NULL,
  `interval` INT NOT NULL,
  until DATETIME,
  count INT
);