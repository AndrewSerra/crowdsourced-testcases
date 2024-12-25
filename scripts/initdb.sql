use csdb;

create table if not exists students (
	id int primary key auto_increment,
    anon_name varchar(255),
    first_name varchar(255) not null,
    last_name varchar(255) not null,
	email varchar(255) not null,
    email_verified tinyint(1) not null default 0,
    created_at datetime not null default now(),

    unique(email)
);

create trigger student_anon_name_set
	before insert on students
	for each row set new.anon_name = uuid_to_bin(uuid());

create table if not exists instructors (
	id int primary key auto_increment,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
	email varchar(255) not null,
    email_verified tinyint(1) not null default 0,
    created_at datetime not null default now()
);

create table if not exists courses (
	id int primary key auto_increment,
    title varchar(255) not null,
    owner_id int not null,
    join_tk varchar(255) not null,
    created_at datetime not null default now(),

    constraint uc_course unique(title, owner_id),
    foreign key (owner_id)
		references instructors(id)
);

create table if not exists course_registration (
	id int primary key auto_increment,
    course_id int not null,
    student_id int not null,
    entry_code binary(16),
    created_at datetime not null default now(),

    constraint uc_registration unique(course_id, student_id),
    foreign key (course_id)
		references courses(id),
	foreign key (student_id)
		references students(id)
);

create trigger register_code_set
	before insert on course_registration
	for each row set new.entry_code = uuid_to_bin(uuid());

create table if not exists assignments (
	id int primary key auto_increment,
    title varchar(255) not null,
    course_id int not null,
    start_date datetime not null,
    end_date datetime not null,
    is_open tinyint(1) not null default 1,
    is_published tinyint(1) not null default 0,
    created_at datetime not null default now(),

    constraint uc_course_assingment unique(title, course_id),
    foreign key (course_id)
		references courses(id)
        on delete cascade
);

create table if not exists assignment_submissions (
	id int primary key auto_increment,
    course_id int not null,
    owner_id int not null,
    assignment_id int not null,
    grading_status tinyint(1) not null default 0,
    submitted_at datetime not null default now(),

    foreign key (course_id)
		references courses(id),
	foreign key (owner_id)
		references students(id),
	foreign key (assignment_id)
		references assignments(id)
);

create table if not exists testcase_submissions (
	id int primary key auto_increment,
    course_id int not null,
    owner_id int not null,
    assignment_id int not null,
    input_data blob not null,
    expected_result blob not null,
    submitted_at datetime not null default now()
);
