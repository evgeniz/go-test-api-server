CREATE TABLE `guru_team`.`users` (
  `id` INT NOT NULL,
  `balance` DECIMAL(6,2) NOT NULL,
  `deposit_count` INT NOT NULL,
  `deposit_sum` DECIMAL(6,2) NOT NULL,
  `bet_count` INT NOT NULL,
  `bet_sum` DECIMAL(6,2) NOT NULL,
  `win_count` INT NOT NULL,
  `win_sum` DECIMAL(6,2) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);
