SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

--
-- Database: `flow-projects`
--

CREATE DATABASE IF NOT EXISTS `flow-projects` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `flow-projects`;

-- --------------------------------------------------------

--
-- Table structure for table `projects`
--

CREATE TABLE `projects` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `theme_color` char(6) NOT NULL,
  `parent_id` bigint UNSIGNED DEFAULT NULL,
  `pinned` tinyint(1) NOT NULL DEFAULT '0',
  `order` tinyint UNSIGNED NOT NULL DEFAULT '0',
  `hidden` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (id)
);