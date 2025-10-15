set sql_safe_updates=0;

-- 数据库创建
create database if not exists course
character set utf8mb4
collate utf8mb4_unicode_ci;


-- 创建用户并设置密码，docker需要配置外界可访问
create user if not exists 'guest'@'%'
identified with mysql_native_password by 'Guest123@';

-- 赋予权限
grant all privileges on course.* to 'guest'@'%';

-- 针对flask db migrate报错
use course;
CREATE TABLE IF NOT EXISTS alembic_version (
    version_num VARCHAR(32) NOT NULL,
    PRIMARY KEY (version_num)
);

-- 刷新
flush privileges;