DROP TABLE classrooms;

create table classrooms (
    week    int  not null,
    weekday int not null,
    building char(1) not null,
    available_classrooms json not null,
    INDEX (week,weekday)
);
