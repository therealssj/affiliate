create table TRACKING_CODE(
	ADDRESS varchar(255) not null,
	ID bigserial NOT NULL unique,
	REF_ADDRESS varchar(255) default null ,
	CREATION timestamp NOT NULL,
	primary key(ADDRESS)
);
alter table TRACKING_CODE add CONSTRAINT  REF_ADDR FOREIGN KEY(REF_ADDRESS) REFERENCES TRACKING_CODE(ADDRESS);

create table OTC_MAPPING(
	ID bigint(20) NOT NULL AUTO_INCREMENT primary key,
	VERSION bigint(20) NOT NULL,
	CREATION timestamp NOT NULL,
	LAST_MODIFIED timestamp NOT NULL,
	ADDRESS varchar(255) not null，
	CODE varchar(32) default null,
	CURRENCY_TYPE varchar(64) not null,
	PAYMENT_ADDR varchar(512) not null,
	PRIVATE_KEY varchar(2048) not null,
	PAY_AMOUNT decimal(9,9) default null,
	BUY_AMOUNT bigint(20) default null,
	REWARD bit(1) NOT NULL default false
)ENGINE=INNODB DEFAULT CHARSET=utf8;

create table REWARD(
	ID bigint(20) NOT NULL AUTO_INCREMENT primary key,
	VERSION bigint(20) NOT NULL,
	CREATION timestamp NOT NULL,
	LAST_MODIFIED timestamp NOT NULL,
	MAPPING_ID bigint(20) not null FOREIGN KEY (MAPPING_ID) REFERENCES OTC_MAPPING(ID),
	AMOUNT bigint(20) not null,
	ADDRESS varchar(255) not null，
	TRANSACTION_ID varchar(255) default null
)ENGINE=INNODB DEFAULT CHARSET=utf8;


