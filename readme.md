
## main
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

```
'4194304'
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
```
直接生成 csv 文件，按行分配，字段使用 , 分割，使用 load file 导入
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

例如：load data infile "/home/mark/Orders txt"replace into table Orders fields terminated by',' enclosed by '"';

（2）lines 关键字指定了每条记录的分隔符默认为'\n'即为换行符

如果两个字段都指定了那fields必须在lines之前。如果不指定fields关键字缺省值与如果你这样写的相同： fields terminated by'\t'enclosed by ’ '' ‘ escaped by'\\'

如果你不指定一个lines子句，缺省值与如果你这样写的相同： lines terminated by'\n'

例如：load data infile "/jiaoben/load.txt" replace into tabletest fields terminated by ',' lines terminated by '/n';


load data infile '/tmp/t0.txt' into table t0 character set gbk fieldsterminated by ',' enclosed by '"' lines terminated by '\n' (`name`,`age`,`description`)set update_time=current_timestamp;
```
对于数据库的链接，就算可以重载 config.toml 也需要重新初始化 database.DB 链接