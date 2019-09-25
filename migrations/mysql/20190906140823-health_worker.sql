
-- +migrate Up
CREATE TABLE health_checks (
  id                    INT AUTO_INCREMENT,
  name                  VARCHAR(255) NOT NULL,
  type                  TINYINT DEFAULT 0,
  check_interval        INT DEFAULT NULL,
  threshould            INT DEFAULT NULL,
  params                JSON DEFAULT NULL,
  PRIMARY KEY (id)
) Engine=InnoDB CHARACTER SET 'latin1';

CREATE UNIQUE INDEX health_checks_name_index ON health_checks(name);

CREATE TABLE routing_policies (
  id                    INT AUTO_INCREMENT,
  record_id             INT NOT NULL,
  health_check_id       INT NOT NULL,
  type                  TINYINT DEFAULT 0,
  PRIMARY KEY (id)
) Engine=InnoDB CHARACTER SET 'latin1';

-- +migrate Down
