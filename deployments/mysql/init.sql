Use api_db;

DROP TABLE IF EXISTS `api_db`.`order_shipping_address`;

CREATE TABLE `api_db`.`order_shipping_address` (
  `order_id` BIGINT NOT NULL,
  `first_name` VARCHAR(45) NOT NULL,
  `address1` VARCHAR(45) NOT NULL,
  `postcode` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`order_id`) 
);

DROP TABLE IF EXISTS `api_db`.`order`;
CREATE TABLE `api_db`.`order` (
  `id` BIGINT NOT NULL,
  `email` VARCHAR(45) NOT NULL,
  `total_price` VARCHAR(45) NOT NULL,
  `total_weight_grams` INT NOT NULL,
  `order_number` BIGINT NOT NULL,
  PRIMARY KEY (`id`)
);

DROP TABLE IF EXISTS `api_db`.`order_delivery`;
CREATE TABLE `api_db`.`order_delivery` (
  `order_id` BIGINT NOT NULL,
  `delivery_id` BIGINT NOT NULL,
  PRIMARY KEY (`order_id`, `delivery_id`)
);

DROP TABLE IF EXISTS `api_db`.`shipping_line`;

CREATE TABLE `api_db`.`shipping_line` (
  `id` BIGINT NOT NULL,
  PRIMARY KEY (`id`));

DROP TABLE IF EXISTS `api_db`.`order_to_shipping_line`;

CREATE TABLE `order_to_shipping_line` (
  `order_id` bigint NOT NULL,
  `shipping_line_id` bigint NOT NULL,
  `title` varchar(45) NOT NULL,
  `price` varchar(45) NOT NULL,
  PRIMARY KEY (`order_id`,`shipping_line_id`)
)