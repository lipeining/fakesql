
## fakesql
### 配置
```go
viper.SetDefault("Xorm", map[string]string{
  "User": "root", 
  "Passwd": "root",
  "Database": "test", 
  "SecurePivFile": "C:/ProgramData/MySQL/MySQL Server 5.7/Uploads",
})
```
链接数据库，使用的是 xorm 包，可以在 config.toml 修改对应的数据库和数据库对应的 securePivFile 路径

### 参数
```go
flag.StringVar(&tblName, "tblName", "user", "tblName")
flag.StringVar(&jsonPath, "jsonPath", "./tables/user.json", "jsonPath absolute or relative path")
flag.StringVar(&num, "num", defaultNum, "generate of num rows")
```
- tblName 创建的表的名字
- jsonPath 表的结构定义 具体参考 /tables 目录，对应的 struct 位于 /moddel
- num 创建的数据行数量

默认命令可以为：(需要将 table 包含在内，避免方法未定义)
go run main.go table.go -tblName=light -jsonPath="./tables/light.json" -num=10000
// 1秒钟
go run main.go table.go -tblName=dark -jsonPath="./tables/dark.json" -num=1000000
// 16秒钟
对应的执行结果
```bash
PS D:\fakesql> go run main.go table.go -tblName=light -jsonPath="./tables/light.json" -num=10000
{fake sql {root root test C:/ProgramData/MySQL/MySQL Server 5.7/Uploads} {root:root@/test?charset=utf8mb4&parseTime=True&loc=Local C:/ProgramData/MySQL/MySQL Server 5.7/Uploads}}
light 10000 ./tables/light.json
./tables/light.json  cols length:  4
CREATE TABLE IF NOT EXISTS test.light (id BIGINT(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "id", name VARCHAR(255) NOT NULL COMMENT "name", power VARCHAR(255) NOT NULL COMMENT "power", create_time datetime NOT NULL COMMENT "create_time") DEFAULT CHARACTER SET utf8mb4
[xorm] [info]  2020/05/07 13:38:55.832807 [SQL] CREATE TABLE IF NOT EXISTS test.light (id BIGINT(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "id", name VARCHAR(255) NOT NULL COMMENT "name", power VARCHAR(255) NOT NULL COMMENT "power", create_time datetime NOT NULL COMMENT "create_time") DEFAULT CHARACTER SET utf8mb4 [] - 886.2325ms
load data infile 'C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/light_10000.txt' replace into table test.light character set utf8mb4 fields terminated by ',' (`id`,`name`,`power`,`create_time`);
[xorm] [info]  2020/05/07 13:38:56.719622 [SQL] load data infile 'C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/light_10000.txt' replace into table test.light character set utf8mb4 fields terminated by ',' (`id`,`name`,`power`,`create_time`); [] - 847.8391ms
load table results []
PS D:\fakesql> go run main.go table.go -tblName=dark -jsonPath="./tables/dark.json" -num=1000000
{fake sql {root root test C:/ProgramData/MySQL/MySQL Server 5.7/Uploads} {root:root@/test?charset=utf8mb4&parseTime=True&loc=Local C:/ProgramData/MySQL/MySQL Server 5.7/Uploads}}
dark 1000000 ./tables/dark.json
./tables/dark.json  cols length:  4
CREATE TABLE IF NOT EXISTS test.dark (id BIGINT(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "id", name VARCHAR(255) NOT NULL COMMENT "name", dark VARCHAR(255) NOT NULL COMMENT "dark", create_time datetime NOT NULL COMMENT "create_time") DEFAULT CHARACTER SET utf8mb4
[xorm] [info]  2020/05/07 13:42:22.752505 [SQL] CREATE TABLE IF NOT EXISTS test.dark (id BIGINT(20) NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "id", name VARCHAR(255) NOT NULL COMMENT "name", dark VARCHAR(255) NOT NULL COMMENT "dark", create_time datetime NOT NULL COMMENT "create_time") DEFAULT CHARACTER SET utf8mb4 [] - 373.5173ms
load data infile 'C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/dark_1000000.txt' replace into table test.dark character set utf8mb4 fields terminated by ',' (`id`,`name`,`dark`,`create_time`);
[xorm] [info]  2020/05/07 13:42:38.595027 [SQL] load data infile 'C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/dark_1000000.txt' replace into table test.dark character set utf8mb4 fields terminated by ',' (`id`,`name`,`dark`,`create_time`); [] - 15.841523s
load table results []
```

### todo
- 制作 cmd 工具
- 制作前端操作页面，支持 http 方式调用
- 针对 gofakeit (github.com/brianvoe/gofakeit/v5) 丰富数据类型
- 考虑 runtine pool 等优化性能方式
- 考虑导出 sql 语句的实现方式
- 测试  insert into values , insert into select 的性能

结合 mysql 的 slow log, explain, sql advisor 等工具，分析 mysql 的 log 日志 bin-log, error-log, normal-log 等，
可以简单实现测试方案的生成。
使用 select * from information.schema 等方式提供数据库的样式，
导出对应的 json 文件。
支持所有的 table column 定义，null not null, unique, indexing, default value, comment,
对应上支持 regrex 的数据生成方式，结合 gofakeit 的内容。
考虑引入已有的 sql advisor 工具的包，生成对应的索引建议，sql改写建议
SOAR	sqlcheck	pt-query-advisor	SQL Advisor	Inception	sqlautoreview

### 理论基础
#### load data
```
代码直接生成 csv 文件，按行分配，字段使用 , 分割，使用 load data 导入
```
```
LOAD DATA [LOW_PRIORITY | CONCURRENT] [LOCAL] INFILE 'file_name'
[REPLACE | IGNORE]
INTO TABLE tbl_name
[PARTITION (partition_name,...)]
[CHARACTER SET charset_name]
[{FIELDS | COLUMNS}
[TERMINATED BY 'string']
[[OPTIONALLY] ENCLOSED BY 'char']
[ESCAPED BY 'char']
]
[LINES
[STARTING BY 'string']
[TERMINATED BY 'string']
]
[IGNORE number {LINES | ROWS}]
[(col_name_or_user_var,...)]
[SET col_name = expr,...]

（1） fields关键字指定了文件记段的分割格式，如果用到这个关键字，MySQL剖析器希望看到至少有下面的一个选项： 
terminatedby分隔符：意思是以什么字符作为分隔符
enclosed by字段括起字符
escaped by转义字符

terminated by描述字段的分隔符，默认情况下是tab字符（\t） 
enclosed by描述的是字段的括起字符。
escaped by描述的转义字符。默认的是反斜杠（backslash：\）  

例如：load data infile "/home/xxx/xxx txt" replace into table Orders fields terminated by',' enclosed by '"';

（2）lines 关键字指定了每条记录的分隔符默认为'\n'即为换行符

如果两个字段都指定了那fields必须在lines之前。如果不指定fields关键字缺省值与如果你这样写的相同： fields terminated by'\t'enclosed by ’ '' ‘ escaped by'\\'

如果你不指定一个lines子句，缺省值与如果你这样写的相同： lines terminated by '\n'

例如：load data infile "/xxx/load.txt" replace into tabletest fields terminated by ',' lines terminated by '/n';

load data infile '/xxx/xxx.txt' into table t0 character set gbk fieldsterminated by ',' enclosed by '"' lines terminated by '\n' (`name`,`age`,`description`) set update_time=current_timestamp;
```
#### insert into select
```sql
use big;
show variables like "%secure%";
load data infile 'C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/base.txt' replace into table tbl_tmp_id;
load data infile 'C:/ProgramData/MySQL/MySQL Server 5.7/Uploads/base.txt' replace into table test.tbl_tmp_id;
select count(id) from tbl_tmp_id;
CREATE TABLE `tbl_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `c_user_id` varchar(36) NOT NULL DEFAULT '',
  `c_name` varchar(22) NOT NULL DEFAULT '',
  `c_province_id` int(11) NOT NULL DEFAULT 0,
  `c_city_id` int(11) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
 INSERT INTO tbl_user
       SELECT
         id,
         uuid(),
         CONCAT('userNickName', id),
         FLOOR(Rand() * 1000),
         FLOOR(Rand() * 100),
         NOW()
       FROM
	 tbl_tmp_id;
 INSERT INTO tbl_user
       SELECT
         id,
         uuid(),
         CONCAT('userNickName', id),
         FLOOR(Rand() * 1000),
         FLOOR(Rand() * 100),
         date_add(create_time, interval FLOOR(1 + (RAND() * 10)) month)
       FROM
	 tbl_tmp_id;
UPDATE tbl_user SET create_time=date_add(create_time, interval FLOOR(1 + (RAND() * 7)) year) where id > 1;
UPDATE tbl_user SET create_time=date_add(create_time, interval FLOOR(1 + (RAND() * 30)) day) where id > 1;
UPDATE tbl_user SET create_time=date_add(create_time, interval FLOOR(1 + (RAND() * 10)) month) where id > 1;

ALTER TABLE `big`.`tbl_user` 
ADD INDEX `create_time_idx` USING BTREE (`create_time`);
explain select * from tbl_user where id < 50;
explain select * from tbl_user where c_province_id > 400;
explain select count(id) from tbl_user where c_province_id > 400;
select count(id) from tbl_user where id < 400;
select * from tbl_user where id = 400;
explain select * from tbl_user where id = 400;
select count(id) from tbl_user where id > 400 and id < 300000;
select count(id) from tbl_user where c_province_id > 400 and c_province_id < 300000;
select count(distinct c_province_id)/count(id), count(distinct c_city_id)/count(id) from tbl_user;
explain select count(distinct c_city_id)/count(id) from tbl_user;
select id from tbl_user where create_time > now();
explain select id from tbl_user where create_time > now();
select count(id) from tbl_user where create_time > now();
select count(id) from tbl_user where create_time > now() limit 1,0;
explain select count(id) from tbl_user where create_time > now();
select create_time, count(id) from tbl_user where create_time > now() group by year(create_time);
explain select create_time, count(id) from tbl_user where create_time > now() group by year(create_time);

substring(MD5(RAND()), 1, 20) 可以获取随机字符串

-- 随机姓名 可根据需要增加/减少样本
set @SURNAME = '王李张刘陈杨黄赵吴周徐孙马朱胡郭何高林罗郑梁谢宋唐位许韩冯邓曹彭曾萧田董潘袁于蒋蔡余杜叶程苏魏吕丁任沈姚卢姜崔钟谭陆汪范金石廖贾夏韦傅方白邹孟熊秦邱江尹薛阎段雷侯龙史陶黎贺顾毛郝龚邵万钱严覃武戴莫孔向汤';
 
set @NAME = '丹举义之乐书乾云亦从代以伟佑俊修健傲儿元光兰冬冰冷凌凝凡凯初力勤千卉半华南博又友同向君听和哲嘉国坚城夏夜天奇奥如妙子存季孤宇安宛宸寒寻尔尧山岚峻巧平幼康建开弘强彤彦彬彭心忆志念怀怜恨惜慕成擎敏文新旋旭昊明易昕映春昱晋晓晗晟景晴智曼朋朗杰松枫柏柔柳格桃梦楷槐正水沛波泽洁洋济浦浩海涛润涵渊源溥濮瀚灵灿炎烟烨然煊煜熙熠玉珊珍理琪琴瑜瑞瑶瑾璞痴皓盼真睿碧磊祥祺秉程立竹笑紫绍经绿群翠翰致航良芙芷苍苑若茂荣莲菡菱萱蓉蓝蕊蕾薇蝶觅访诚语谷豪赋超越轩辉达远邃醉金鑫锦问雁雅雨雪霖霜露青靖静风飞香驰骞高鸿鹏鹤黎';
 
-- length(@surname)/3 是因为中文字符占用3个长度
select concat(substr(@surname,floor(rand()*length(@surname)/3+1),1), substr(@NAME,floor(rand()*length(@NAME)/3+1),1), substr(@NAME,floor(rand()*length(@NAME)/3+1),1));

```
#### insert into values 
在客户端直接生成大数据插入语句，发送到 mysql 服务器进行执行
```
show variables like "%max_allowed_packet%";
max_allowed_packet = '4194304' 
通过设置数据包大小，控制批量插入数据的大小数量 64MB 以上
max_allowed_packet 数据包大小
innodb_log_buffer_size 事务大小
这样的话，可以使用 fake 数据之类的工具，生成一个 
START TRANSACTION;
insert into tabl (`id`,`xx`, ``,) values(),(),();
commit;
的大 sql 批量插入语句，然后执行即可。
好处在于：数据格式可以定制。
```
#### mysqldump mysqlimport 
生成固定格式的 sql 文件，能够通过外部工具进行读写，导入导出。