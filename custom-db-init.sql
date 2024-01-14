create table listChannels (
    id uuid primary key,
    channel_name varchar(255),
    public bool
);

create table users (
    id uuid primary key,
    username varchar(255),
    cookie varchar(1023),
    online bool,
    lastChannel uuid,
    constraint fk_channel foreign key (lastChannel) references listChannels (id) on update cascade
);

create table allowedChannel (
    id serial primary key,
    channelID uuid,
    userID uuid,
    constraint fk_user foreign key (userID) references users (id) on delete cascade on update cascade,
    constraint fk_channel foreign key (channelID) references listChannels (id) on delete cascade on update cascade
);

alter table
    users
add
    unique (username);

alter table
    listChannels
add
    unique (channel_name);

insert into
    listChannels (id, channel_name, public)
VALUES
    (gen_random_uuid(), 'public_1', true);

insert into
    listChannels (id, channel_name, public)
VALUES
    (gen_random_uuid(), 'general', true);

insert into
    listChannels (id, channel_name, public)
VALUES
    (gen_random_uuid(), 'partychat', false);

insert into
    listChannels (id, channel_name, public)
VALUES
    (gen_random_uuid(), 'private_1', false);

create table public_1 (
    id serial primary key,
    username varchar(255),
    stamp timestamp DEFAULT CURRENT_TIMESTAMP,
    msg text,
    constraint fk_user foreign key (username) references users (username) on delete cascade on update cascade
);

create table general (
    id serial primary key,
    username varchar(255),
    stamp timestamp DEFAULT CURRENT_TIMESTAMP,
    msg text,
    constraint fk_user foreign key (username) references users (username) on delete cascade on update cascade
);

create table private_1 (
    id serial primary key,
    username varchar(255),
    stamp timestamp DEFAULT CURRENT_TIMESTAMP,
    msg text,
    constraint fk_user foreign key (username) references users (username) on delete cascade on update cascade
);

create table partychat (
    id serial primary key,
    username varchar(255),
    stamp timestamp DEFAULT CURRENT_TIMESTAMP,
    msg text,
    constraint fk_user foreign key (username) references users (username) on delete cascade on update cascade
);

insert into users (id, username) values (gen_random_uuid(), 'admin');
insert into users (id, username) values (gen_random_uuid(), 'userparty');
insert into users (id, username) values (gen_random_uuid(), 'usercool');

insert into allowedChannel (channelID, userID) values ((select id from listChannels where channel_name='partychat'), (select id from users where username='admin'));
insert into allowedChannel (channelID, userID) values ((select id from listChannels where channel_name='private_1'), (select id from users where username='admin'));
insert into allowedChannel (channelID, userID) values ((select id from listChannels where channel_name='partychat'), (select id from users where username='userparty'));
insert into allowedChannel (channelID, userID) values ((select id from listChannels where channel_name='private_1'), (select id from users where username='usercool'));