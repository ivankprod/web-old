SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

CREATE DATABASE IF NOT EXISTS `ivankprodru` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `ivankprodru`;

CREATE TABLE `users` (
  `user_id` int(11) UNSIGNED NOT NULL,
  `user_group` int(11) NOT NULL,
  `user_social_id` varchar(255) NOT NULL,
  `user_access_token` text NOT NULL,
  `user_avatar_path` text NOT NULL,
  `user_email` varchar(255) NOT NULL,
  `user_name_first` varchar(255) NOT NULL,
  `user_name_last` varchar(255) NOT NULL,
  `user_last_access` varchar(255) NOT NULL,
  `user_type` int(11) NOT NULL,
  `user_role` int(11) NOT NULL DEFAULT 2
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Users table';

INSERT IGNORE INTO `users` (`user_id`, `user_group`, `user_social_id`, `user_access_token`, `user_avatar_path`, `user_email`, `user_name_first`, `user_name_last`, `user_last_access`, `user_type`, `user_role`) VALUES
(2, 2, '151552597', '147db9f6b024d1aaae5ea848313ae6febd0accaf73a29d87f8007446499133303694002b8220f0bc054db', 'https://sun1-84.userapi.com/s/v1/ig2/HhkKEHvxEcJEcvuOun3WnFaG-OkKhwuOnCJxnXMpuNmUGm0ebtPH4_OAnDQ18wMz0b9vxLYZ6vEbQgy_bKCCYViV.jpg?size=200x0&quality=96&crop=226,1,719,719&ava=1', 'Tu134music@mail.ru', 'Алексей', 'Егоров', '2021-04-02 21:43:54', 0, 2),
(10, 10, '53625544', '163ff30b2c8a65677a87ed15ea9f12a5829a2e3723aab15691f4676cdd9187bb6fbedbe018f82dc0fd14e', 'https://sun1-89.userapi.com/s/v1/ig2/lg7PtFBJx9X_THTvDrctbjFWCKfsvwek5qkxBwdXTFqh2Ova7RHPRuZlCtOJIkdDzos5QoDH3Lup65ZKQIGW04bY.jpg?size=400x0&quality=96&crop=549,0,722,1080&ava=1', 'houseprotector@mail.ru', 'Иван', 'Кулаков', '2021-07-26 03:02:05', 0, 3),
(13, 10, '100347571556835054591', 'ya29.a0ARrdaM9GRPyI5DtgvNkxJK0DQqxxQOI4y5bY0vOi48aGYdreOcxTblu3V_g2i7v_aCqWo5FDRVULf7XdzJ4CCxt0WdnjcOt6Pz0nO9Mf_ZQfiZ1IC2Y04KV3Dcn8_KS5IyXfoUE4AsPEzqMqCPGazZAzhtF2ezc', 'https://lh3.googleusercontent.com/a-/AOh14GggbjyduqSJCKTD2L1-GJOy-zmszQ0GSB5Du_Yc3w=s400-c', 'colldierofficial@gmail.com', 'IvanK', 'Production', '2021-07-26 22:06:55', 3, 3),
(14, 10, '2487408934736154', 'EAAOUe8LI7CsBAMfZBVXZCD5ZBZAv2v1CqDEltnmvg8EIs6DZA9jbMmAolx0JCxvVVBDO7cxaKGDZC4n7XxIPddUfCuUcDGaZCerdu7MEfHhoA28axFZAGiAjo2biGEVsT8MJEVL8UwrxoUZCpWUxXwqeEd9Aatf5BcmOwdD3Vf8mtFAZDZD', 'https://platform-lookaside.fbsbx.com/platform/profilepic/?asid=2487408934736154&width=400&ext=1629761757&hash=AeQhxGiv6qrRuyXXeP8', 'houseprotector@mail.ru', 'Иван', 'Кулаков', '2021-07-25 02:40:28', 2, 3),
(15, 10, '380639006', 'AQAAAAAWsBceAAdFE4JYO3hkj0JFjmU9ZsdUcp0', 'https://avatars.yandex.net/get-yapic/35885/PfSisKFz5rzvCnCp6UVlxHwY-1/islands-200', 'ivankprod@yandex.ru', 'Иван', 'Кулаков', '2021-07-27 05:14:32', 1, 3);

CREATE TABLE `users_roles` (
  `id` int(11) NOT NULL,
  `role` text NOT NULL,
  `sort` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Users roles table';

INSERT IGNORE INTO `users_roles` (`id`, `role`, `sort`) VALUES
(1, 'Заблокирован', 4),
(2, 'Гость', 3),
(3, 'Веб-мастер', 1),
(4, 'Администратор', 2);

CREATE TABLE `users_types` (
  `id` int(11) NOT NULL,
  `type` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Users types table';

INSERT IGNORE INTO `users_types` (`id`, `type`) VALUES
(0, 'ВКонтакте'),
(1, 'Яндекс'),
(2, 'Facebook'),
(3, 'Google');


ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`) USING BTREE,
  ADD KEY `user_type` (`user_type`) USING BTREE,
  ADD KEY `user_role` (`user_role`);

ALTER TABLE `users_roles`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `users_types`
  ADD PRIMARY KEY (`id`);


ALTER TABLE `users`
  MODIFY `user_id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

ALTER TABLE `users_types`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;


ALTER TABLE `users`
  ADD CONSTRAINT `users_ibfk_1` FOREIGN KEY (`user_type`) REFERENCES `users_types` (`id`),
  ADD CONSTRAINT `users_ifbk2` FOREIGN KEY (`user_role`) REFERENCES `users_roles` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
