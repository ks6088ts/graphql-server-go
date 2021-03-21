create table company
(
    company_cd integer not null,
    rr_cd integer not null,
    company_name varchar not null,
    company_name_k varchar not null,
    company_name_h varchar not null,
    company_name_r varchar not null,
    company_url varchar,
    company_type integer,
    e_status integer,
    e_sort integer,
    PRIMARY KEY (company_cd)
);
comment on table company is 'company.csv';

create table station_join
(
    line_cd integer not null,
    station_cd1 integer not null,
    station_cd2 integer not null
);
comment on table station_join is 'join.csv';

create table line
(
    line_cd integer not null,
    company_cd integer not null,
    line_name varchar not null,
    line_name_k varchar not null,
    line_name_h varchar not null,
    line_color_c varchar,
    line_color_t varchar,
    line_type varchar,
    lon float not null,
    lat float not null,
    zoom integer,
    e_status integer,
    e_sort integer,
    PRIMARY KEY (line_cd)
);
comment on table line is 'line.csv';

create table station
(
    station_cd integer not null,
    station_g_cd integer not null,
    station_name varchar not null,
    station_name_k varchar,
    station_name_r varchar,
    line_cd integer,
    pref_cd integer,
    post varchar,
    address varchar,
    lon float,
    lat float,
    open_ymd varchar ,
    close_ymd varchar,
    e_status integer,
    e_sort integer,
    PRIMARY KEY (station_cd)
);
comment on table station is 'station.csv';