create database db_tars;
create database db_user_system;
create database db_tars_web;
create database db_cache_web;
create database tars_stat;
create database tars_property;

use db_tars;
source ../deploy/deploy/framework/sql/db_tars.sql;
use db_user_system;
source ../deploy/deploy/web/demo/sql/db_user_system.sql;
use db_tars_web;
source ../deploy/deploy/web/sql/db_tars_web.sql;
use db_cache_web;
source ../deploy/deploy/web/sql/db_cache_web.sql;
use tars_stat;
source ../deploy/deploy/web/sql/db_cache_web.sql;
