create table listChannels
(
    id           uuid primary key,
    channel_name varchar(255),
    public       bool
);

create table users
(
    id          uuid primary key,
    username    varchar(255),
    cookie      varchar(255),
    online      bool,
    lastChannel uuid,
    constraint fk_channel foreign key (lastChannel) references listChannels (id) on update cascade
);

create table allowedChannel
(
    id        serial primary key,
    channelID uuid,
    userID    uuid,
    constraint fk_user foreign key (userID) references users (id) on delete cascade on update cascade,
    constraint fk_channel foreign key (channelID) references listChannels (id) on delete cascade on update cascade
);

alter table users
    add unique (username);

alter table listChannels
    add unique (channel_name);

insert into listChannels (id, channel_name, public)
VALUES (gen_random_uuid(), 'public_1', true);
insert into listChannels (id, channel_name, public)
VALUES (gen_random_uuid(), 'private_1', false);

create table public_1
(
    id        serial primary key,
    username  varchar(255),
    timestamp timestamp DEFAULT CURRENT_TIMESTAMP,
    message   text,
    constraint fk_user foreign key (username) references users (id) on delete cascade on update cascade,
);

create table private_1
(
    id        serial primary key,
    username  varchar(255),
    timestamp timestamp DEFAULT CURRENT_TIMESTAMP,
    message   text,
    constraint fk_user foreign key (username) references users (id) on delete cascade on update cascade,
);