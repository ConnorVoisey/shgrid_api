-- This needs to be moved into something that can run migrations
create table addresses
(
    id             uuid         not null
        constraint addresses_pk
            primary key,
    address_line_1 varchar(255),
    address_line_2 varchar(255),
    town_or_city   varchar(255) not null,
    country        varchar(255) not null,
    postcode       varchar(100) not null
);

create table organisations
(
    id         uuid not null
        constraint organisations_pk
            primary key,
    name       text not null,
    address_id uuid
        constraint organisations_addresses_id_fk
            references addresses
);

create table contacts
(
    id              uuid         not null
        constraint contacts_pk
            primary key,
    first_name      text,
    last_name       text         not null,
    email           varchar(255) not null,
    organisation_id uuid
        constraint contacts_organisations_id_fk
            references organisations,
    address_id      uuid
        constraint contacts_addresses_id_fk
            references addresses
);

create index contacts_email_index
    on contacts (email);

create index organisations_name_index
    on organisations (name);


