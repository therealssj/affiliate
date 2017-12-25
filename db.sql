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
	SEQ bigint not null,
	UPDATED_AT bigint not null,
	TRANSACTION_ID varchar(512) not null,
	DEPOSIT_AMOUNT NUMERIC(18,9) not null,
	BUY_AMOUNT bigint not null,
	RATE NUMERIC(18,9) not null,
	HEIGHT bigint not null,
	primary key(ID)
)
alter table DEPOSIT_RECORD add CONSTRAINT  MAPPING_ID FOREIGN KEY(MAPPING_ID) REFERENCES BUY_ADDR_MAPPING(ID);

create table AWARD_RECORD(
	ID bigserial NOT NULL,
	SEND_TIME timestamp NOT NULL,
	DEPOSIT_ID bigint not null,
	ADDRESS varchar(255) not null,
	AMOUNT bigint not null,
	REWARD boolean NOT NULL default false,
	REWARD_TYPE varchar(32) default null,
	TRANSACTION_ID varchar(255) default null,
	primary key(ID)
);
alter table SEND_COIN_RECORD add CONSTRAINT  DEPOSIT_ID FOREIGN KEY(DEPOSIT_ID) REFERENCES DEPOSIT_RECORD(ID);

create table KV_STORE(
	NAME varchar(64) not null,
	INT_VAL bigint default null,
	STR_VAL varchar(255) default null,
	primary key(NAME)
);


