create table TRACKING_CODE(
	ADDRESS varchar(255) not null,
	ID bigserial NOT NULL unique,
	REF_ADDRESS varchar(255) default null ,
	CREATION timestamp NOT NULL,
	primary key(ADDRESS)
);
alter table TRACKING_CODE add CONSTRAINT  REF_ADDR FOREIGN KEY(REF_ADDRESS) REFERENCES TRACKING_CODE(ADDRESS);

create table ALL_CRYPTOCURRENCY(
	SHORT_NAME varchar(32) not null,
	FULL_NAME varchar(128) not null,
	RATE NUMERIC(18,9) not null,
	primary key(SHORT_NAME)
);

create table PREPARATIVE_ADDR(
	CURRENCY_TYPE varchar(32) not null,
	DEPOSIT_ADDR varchar(512) not null,
	PRIMARY KEY (CURRENCY_TYPE, DEPOSIT_ADDR)
);
alter table PREPARATIVE_ADDR add CONSTRAINT  CURRENCY_TYPE FOREIGN KEY(CURRENCY_TYPE) REFERENCES ALL_CRYPTOCURRENCY(SHORT_NAME);

create table BUY_ADDR_MAPPING(
	ID bigserial NOT NULL,
	VERSION bigint NOT NULL,
	CREATION timestamp NOT NULL,
	LAST_MODIFIED timestamp NOT NULL,
	ADDRESS varchar(255) not null,
	CURRENCY_TYPE varchar(32) not null,
	DEPOSIT_ADDR varchar(512) not null,
	REF varchar(32) default null,
	DEPOSIT_AMOUNT NUMERIC(18,9) default null,
	BUY_AMOUNT bigint default null,
	LAST_UPDATED timestamp default null,
	TRANSACTION_IDS varchar(512) default null,
	SENT_COIN boolean NOT NULL default false,
	primary key(ID),
	UNIQUE (ADDRESS, CURRENCY_TYPE)
);
alter table BUY_ADDR_MAPPING add CONSTRAINT  CURRENCY_TYPE FOREIGN KEY(CURRENCY_TYPE) REFERENCES ALL_CRYPTOCURRENCY(SHORT_NAME);

create table SEND_COIN_RECORD(
	ID bigserial NOT NULL,
	SEND_TIME timestamp NOT NULL,
	MAPPING_ID bigint not null,
	ADDRESS varchar(255) not null,
	AMOUNT bigint not null,
	REWARD boolean NOT NULL default false,
	REWARD_TYPE varchar(32) default null,
	TRANSACTION_ID varchar(255) default null,
	primary key(ID)
);
alter table SEND_COIN_RECORD add CONSTRAINT  MAPPING_ID FOREIGN KEY(MAPPING_ID) REFERENCES BUY_ADDR_MAPPING(ID);

create table KV_STORE(
	NAME varchar(64) not null,
	INT_VAL bigint default null,
	STR_VAL varchar(255) default null,
	primary key(NAME)
);


