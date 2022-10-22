CREATE TABLE vcs_server(
    server_id serial,
	server_name varchar(50),
	status bool,
	created_time int,
	last_updated int ,
	ipv4 varchar(50),
	PRIMARY KEY (server_id)
);