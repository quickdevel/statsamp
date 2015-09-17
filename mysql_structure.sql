CREATE TABLE IF NOT EXISTS `history` (
  `history_id` int(11) NOT NULL AUTO_INCREMENT,
  `history_servers_total` int(11) DEFAULT NULL,
  `history_servers_online` int(11) DEFAULT NULL,
  `history_slots_total` int(11) DEFAULT NULL,
  `history_slots_used` int(11) DEFAULT NULL,
  `history_date` datetime DEFAULT NULL,
  PRIMARY KEY (`history_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `servers` (
  `server_id` int(11) NOT NULL AUTO_INCREMENT,
  `server_addr` varchar(32) DEFAULT NULL,
  `server_name` varchar(64) DEFAULT NULL,
  `server_players` int(11) DEFAULT NULL,
  `server_maxplayers` int(11) DEFAULT NULL,
  PRIMARY KEY (`server_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
