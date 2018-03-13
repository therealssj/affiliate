create table TRACKING_CODE(
	ADDRESS varchar(255) not null,
	ID bigserial NOT NULL unique,
	REF_ADDR varchar(255) default null ,
	CREATION timestamp NOT NULL,
	primary key(ADDRESS)
);
alter table TRACKING_CODE add CONSTRAINT  REF_ADDR FOREIGN KEY(REF_ADDR) REFERENCES TRACKING_CODE(ADDRESS);

create table ALL_CRYPTOCURRENCY(
	SHORT_NAME varchar(32) not null,
	FULL_NAME varchar(128) not null,
	RATE varchar(64) not null,
	UNIT_POWER smallint not null,
	ENABLED boolean not null,
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
	CREATION timestamp NOT NULL,
	ADDRESS varchar(255) not null,
	CURRENCY_TYPE varchar(32) not null,
	DEPOSIT_ADDR varchar(512) not null,
	REF varchar(32) default null,
	primary key(ID),
	UNIQUE (ADDRESS, CURRENCY_TYPE)
);
alter table BUY_ADDR_MAPPING add CONSTRAINT  CURRENCY_TYPE FOREIGN KEY(CURRENCY_TYPE) REFERENCES ALL_CRYPTOCURRENCY(SHORT_NAME);

create table DEPOSIT_RECORD(
	ID bigserial NOT NULL,
	CREATION timestamp NOT NULL,
	MAPPING_ID bigint not null,
	BUY_ADDR varchar(255) not null,
	CURRENCY_TYPE varchar(32) not null,
	REF_ADDR varchar(255) not null,
	SUPERIOR_REF_ADDR varchar(255) not null,
	SEQ bigint not null unique,
	UPDATED_AT bigint not null,
	TRANSACTION_ID varchar(512) not null,
	DEPOSIT_AMOUNT bigint not null,
	BUY_AMOUNT bigint not null,
	RATE varchar(64) not null,
	HEIGHT bigint not null,
	primary key(ID),
	UNIQUE (MAPPING_ID, TRANSACTION_ID)
);
alter table DEPOSIT_RECORD add CONSTRAINT  CURRENCY_TYPE FOREIGN KEY(CURRENCY_TYPE) REFERENCES ALL_CRYPTOCURRENCY(SHORT_NAME);
alter table DEPOSIT_RECORD add CONSTRAINT  MAPPING_ID FOREIGN KEY(MAPPING_ID) REFERENCES BUY_ADDR_MAPPING(ID);
CREATE INDEX BUY_ADDR ON DEPOSIT_RECORD (BUY_ADDR);
CREATE INDEX CURRENCY_TYPE ON DEPOSIT_RECORD (CURRENCY_TYPE);
CREATE INDEX REF_ADDR ON DEPOSIT_RECORD (REF_ADDR);
CREATE INDEX SUPERIOR_REF_ADDR ON DEPOSIT_RECORD (SUPERIOR_REF_ADDR);

create table REWARD_RECORD(
	ID bigserial NOT NULL,
	CREATION timestamp NOT NULL,
	DEPOSIT_ID bigint not null,
	ADDRESS varchar(255) not null,
	CAL_AMOUNT bigint not null,
	SENT_AMOUNT bigint not NULL,
	SENT_TIME timestamp default NULL,
	SENT boolean not null,
	REWARD_TYPE varchar(32) not null,
	CHECKSUM varchar(255) default null,
	primary key(ID),
	UNIQUE (DEPOSIT_ID, ADDRESS, REWARD_TYPE)
);
alter table REWARD_RECORD add CONSTRAINT  DEPOSIT_ID FOREIGN KEY(DEPOSIT_ID) REFERENCES DEPOSIT_RECORD(ID);
CREATE INDEX ADDRESS ON REWARD_RECORD (ADDRESS);
CREATE INDEX SENT ON REWARD_RECORD (SENT);
CREATE INDEX SENT_AMOUNT ON REWARD_RECORD (SENT_AMOUNT);

create table REWARD_REMAIN(
	ADDRESS varchar(255) not null,
	CREATION timestamp NOT NULL,
	LAST_MODIFIED timestamp NOT NULL,
	AMOUNT bigint not null,
	CHECKSUM varchar(255) not null,
	primary key(ADDRESS)
);

create table KV_STORE(
	NAME varchar(64) not null,
	INT_VAL bigint default null,
	STR_VAL varchar(255) default null,
	primary key(NAME)
);

create table NEWSLETTER(
	EMAIL varchar(256) not null,
	CONCERN_MINER boolean not null,
	CREATION timestamp NOT NULL,
	primary key(EMAIL)
);



