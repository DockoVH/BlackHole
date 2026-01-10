CREATE KEYSPACE IF NOT EXIST blackhole 
WITH REPLICATION {
	'class': 'SimpleStrategy',
	'replication_factor': 1
};

USE blackhole;

CREATE TABLE IF NOT EXISTS igraci {
	igrac_uuid UUID PRIMARY KEY,
	
}
