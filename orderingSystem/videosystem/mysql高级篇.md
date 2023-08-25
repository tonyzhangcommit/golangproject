远程登录连接失败的问题汇总
    首先测试与远程服务器网络是否连接
    测试远程服务器指定端口号是否开放
        windows 电脑上使用telnet 命令进行测试
            telnet remote_ip_address 3306
    关闭windows 防火墙
    关闭远程服务器防火墙  或者 开放指定端口
        sudo firewall-cmd --zone=public --add-port=3306/tcp --permanent
        sudo firewall-cmd --reload
        sudo：以管理员权限执行命令。
        firewall-cmd：用于管理防火墙规则的命令。
        --zone=public：指定防火墙区域，这里使用公共区域。
        --add-port=3306/tcp：添加一个规则，允许 TCP 协议的 3306 端口。
        --permanent：将规则持久化，使其在防火墙重启后仍然有效。
    mysql 默认不允许root 远程登录，只允许登录本机
        update user set host = '%' where user = 'root' 生产环境注意配置ip白名单
        flush privileges;    刷新权限
    8.0 默认对密码进行加密操作，如果不设置会导致远程登录失败、
    alert user 'root'@'%' identified with mysql_native_password by 'PWD'
字符集的修改和底层原理
    MYSQL 默认字符集  latin1(5.7版本默认，需要更改，更改方法，修改my.cnf,然后重启) utf8
    有不同级别的字符集
    服务器级别
    数据库级别
    表级别
    列级别
    utf8
    utf8mb3  utf8别名（mysql）阉割版的utf8, 只使用1-3个字节表示字符
    utf8mb4  正宗utf8，可以存放emoje表情等 1-4
    还要注意字符集的比较规则，有多种，需要了解下
文件目录
    find / -name mysql
    数据库文件保存路径  /var/lib/mysql
    数据库表存储硬盘方式：
        一个表结构  .frm
        一个数据    .ibd(5.7)   ibdtablename (8.0)  系统表空间，8.0将表结构和表空间合并为一个
windows  vs  linux
    windows 大小写不敏感   linux 数据库名，表名，字段名 大小写敏感

用户管理：
    登录：
        mysql -hhostname -uuser -Pport -p
        select * from users\G 纵向显示结果
    创建用户：
        use mysql;
        create user 'username' identified by 'password'   # 新建的用户没有任何权限，注意user表主键为host+user
    更新用户：
        常规update语句
        flush privileges;
    删除用户：
        drop user 'username'@'hostname' @后面可选
    修改密码：
        修改自己的密码
            alter user user() identified by 'newpwd'  方法一
            set password='newpwd'  方法二
        修改其他用户的密码（高权限操作低权限）
            set password for 'user'@'host'='newpwd'
        设置密码过期
            alter USER 'user'@'host' password expire
            设置密码过期，可以登录但是无法进行查询，需要重设密码
    权限管理
        授权两种方式
            直接赋权
            通过赋予角色
            查看权限
                show grants;
            赋予权限
                grant select,update...on database.table to 'username'@'host'
                grant all privileges on *.* to 'username'@'host'
            回收权限
                revoke select，update... on database.table from 'username'@'host'
    角色（权限的集合）
        创建角色
            create role rolename 'rolename'@'host'
        给角色赋予权限
            grant privileges on database.table to 'rolebname'@'host'
        查看角色权限
            show grants for 'rolename'@'host'
        回收角色权限
            revoke select，update... on database.table from 'rolename'@'host'
        删除角色
            drop role 'rolename'@'host'
        赋予用户角色
            grant rolename1,rolename2,...to username1,username2
        查看用户角色是否成功
            show grants for username
        激活角色
            查看角色是否激活成功
                select current_role()
            激活角色
                set default role 'rolename'@'host' to 'username'@'host'  需要退出登录才能生效
            激活角色（另一种方式）
                set global activate_all_roles_on_login=ON   角色永诀激活
        撤销角色
            revoke 'rolename'@'host' from 'username'@'host'  需要高一级的用户才能执行
配置文件
    暂定
sql 语句的执行流程
    1.查询缓存(大部分缓存功能比较鸡肋，它需要严格匹配查询语句，并且涉及系统函数时并不被缓存，再有就有缓存失效的问题，比如表结构修改，涉及updata insert等)，一般建议在静态表中使用缓存，默认不开启查询缓存，可以设置按需缓存，query_cache_type 设置成demand, 代表查询语句中显示的写有sql_cache 关键字时才进行缓存（8.0取消）
        查看mysql 是否开启了查询缓存
            show variables like '%query_cache_type';
    2.解析器
        词法分析
            分析语句中字符串分别是什么，代表什么
        语法分析
            经过词法分析器后，假如语句中每个词都正确，这时根据语法规则，判断输入的语句是否满足mysql 语法,语法通过后，会生成对应的查询语法树
    3.查询优化器
        生成查询树后，在优化器中，决定sql执行路径，比如说是根据全表检索，索引检索等，作用是找到最好的执行计划
    4.执行计划
        根据上步生成的执行计划，在执行之前首先判断权限等，选择存储引擎
    5.查询执行引擎
    6.调用API
打开sql执行记录各个环节记录设置
    查看当前是否开启  select @@session.profiling;
    打开   set profiling=1
    查看所有执行记录
        show profiles;
    查看某一条：
        show profile {cpu,io 等字段可选} for profile_id;
存储引擎
    存储引擎，本质上就是指表的类型，以前叫做表处理器，功能就是接受上层的命令执行相应的功能
    查看mysql 支持什么引擎
        show engines;
        innoDB (支持外键，支持事务，仅仅InnoDB 支持，优先使用)
        相比较与MyISAM（适用于读写场景的需求）引擎来讲， InnoDB 写的效率相对差一些且对内存要求较高，并且会占用更多的磁盘空间以保存数据和索引
    支持对数据库，表设定特定的引擎

索引和调优（重点!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!）
索引数据结构，本质上是一种数据结构。
    索引是在存储引擎中实现的，存储索引可以定义每个表的最大索引数和最大索引长度
优点：
    提高检索效率，降低数据库IO成本（耗时的主要原因）
    通过创建唯一性索引，可以保证数据库表的每一行的数据唯一性（实现此功能还有另外一种方式：定义表时使用unique关键字进行约束）
    对于有依赖关系的子表和父表联合查询时，可以提高查询效率
    使用分组和排序字句进行查询时，可以显著减少查询中分组和排序时间
缺点
    创建和维护索引都需要耗费时间，随数据量的增加，耗费时间也增加
    索引占据磁盘空间
    索引可以提高查询速度，但是会降低更新表的速度（有一个技巧，在大量插入的情境下， 可以先删除表的索引，等写入完成再重新创建索引）
相关概念：
    数据页
    记录
    目录项/数据页 以及随着数据页的增加，目录页也会随之增加，类比于套娃模式增加
    面试点：
        为什么说B+数最高为四层？
            假设一个数据页大小为16k,可以容纳100条数据， 目录页可以存1000条数据，则：
                只有一层：最多100条
                两层：最多1000 * 10
                三层：最多 1000 * 1000 *100
                四层：1000 * 1000 *1000 *100
    索引按照物理实现方式：
        聚簇
            并不是一种单独的索引类型，而是一种数据存储方式。用户记录都在叶子节点，也就是索引即数据，数据即索引
            术语：聚簇，表示数据行和相邻的键值聚簇的存储在一起
            特点：
                使用记录主键值的大小进行记录和页的排序
                1.页内记录按照主键大小顺序排成一个单项链表
                2.各个存放用户记录页也是根据用户记录的主键大小顺序排成一个双向链表
                3.存放目录项记录的页，分为不同层次。同一层次中的页也是根据页中的目录项主键的大小顺序排成一个双向链表
                4.整个过程是InnoDB 自动完成，基于主键，不需要显示声明创建
            优点：
                数据访问更快，因为聚簇索引将索引和数据保存在同一个B+树中， 因此从聚簇索引中获取数据比非聚簇索引更快
                聚簇索引基于主键的排序和查找速度非常快
                由于聚簇索引的特点，在查询一定范围内的数据，由于数据紧密相连，数据库不用从多个数据块中提取，节省大量IO操作
            缺点：
                插入速度严重依赖插入顺序。如果插入导致页分裂，会严重影响性能。所以InnoDB 一般设置主键为递增类型
                更新主键代价很高，一般定义为主键不可更新
                二级索引访问需要两次索引查找，第一次查主键值，第二次查行数据
            限制：
                只有InnoDB 支持

        非聚簇（二级索引，辅助索引）
            区别聚簇索引，叶子节点存放的不是数据，而是选中的列的信息加上主键，执行二级索引时，需要进行一次回表操作（也就是获取ID后重新查询一次聚簇索引）。
        联合索引：
            同时为多个列建立索引，
        hash 索引
        innoDB 自带自适应hash索引功能
        二叉搜索树：
            特性：
                左节点比根节点小，右节点比根节点大，一个根节点最多有两个子节点
            弊端：
                极限情况，退化为链表，为了优化，需要降低树的高度。引入平衡二叉树AVL树
        AVL 树（平衡二叉树：是个空树或者左右两个子树的高度绝对值不超过1，并且左右两个子树都是平衡二叉树（平衡二叉搜索树，红黑树，数堆，伸展树））：
        
        B-tree
            平衡树：
                一个节点最多包含M个子节点，M为B树的阶
                叶子节点和根节点都存放数据
        InnoDB数据存储结构：
            数据页是磁盘和内存互换的基本单位,数据库IO操作最小单位是页（一个页中可以存多行数据），本质上是减少IO操作
            数据页之间双向链表连接，页中记录是通过单链表进行连接，并且根据主键值进行从小到大排序
            InnoDB 将数据划分为若干个页，默认页大小为16K: show varibales like '%innodb_page_size%'
            区是页的上级单位，一个区分配64个连续的页，大小也就是1MB,
            段是区的上层结构，区在内存中是一个连续分配空间，innoDB 默认为64个页，不要求区是连续结构。
            段上级是表空间
        页内部结构:
            文件头
            文件尾
            行格式的记录头信息
        索引分类：
            普通索引
            唯一索引
            主键索引
            全文索引
        创建索引：
            创建表时创建索引
                隐式创建：
                    字段中声明主键，唯一性约束时会自动添加相关列的索引
                显示创建：
                    创建表最后加上 {unique|}index index_name(fields1,|fields2)
                查看 show index from table_name;
                唯一约束自动生成唯一索引，添加唯一索引就有唯一约束
                主键索引，通过删除主键约束可以删除主键索引
                联合索引，可能某种情况导致索引无效（最左前置原则，原因是聚合索引B+树的实现过程）
                explain + sql 语句，性能分析工具
            ALTER TABLE 在已存在的表中创建索引
                命令：
                    alter table_name add {unique 这里可以加索引类型} index index_name(fields1,fields2....) 
            create index 在已存在表中创建索引
                命令：
                    create {unique 索引类型} index index_name on table_name(field1,field2)
        删除索引：
            优化的一种是，在某些有大量更改表的操作时，可以先删除索引，在数据操作完成后再建立索引，提高效率
            show index from table_name;  查看当前表存在的索引
            删除方式1：
                alter table table_name drop index index_name;  有auto_increment 约束的唯一索引不能删除
            方式2：
                drop index index_name on table_name;
            特殊情况，假如存在联合索引，存在三个字段，想要删除索引中某个字段，可以使用以下命令：
                alter table table_name drop column fields_name; 来实现，假如剩一个字段，则联合索引退化为单列索引
        特殊索引：
            降序索引  （8.0新特性）
                应用场景，配置联合索引时，可以显示的设置字段一为升序（默认asc），字段二为降序的索引desc
        隐藏索引：新特性
            为了系统安全，防止删除索引造成出现错误，同时还可以检查添加索引后的性能
            创建索引时，最后加INvisible 关键字 ，默认为visible 可见的
        修改隐藏属性：
            alter table tablename alter index indexname in visible / invisible;
            优化一方面：删除掉长期不用的隐藏索引。因为隐藏索引在数据更新时也在动态更新
    索引设计原则（重点。。。。。。。。。。。。。。。。。）
        适合加索引的字段：
            1.字段性质中有唯一性的限制（约束），alibaba 的设计原则：业务上具有唯一特性的字段，即使是组合字段，也必须建立成唯一索引，不要以为唯一索引影响了insert速度，这种损耗可以忽略不计
            2.频繁作为where查询条件的字段，性能提升一般几十倍
            3.经常order by group by 的字段，这两个函数执行时， mysql 首先将关联字段满足的数据筛选出来然后进行后续操作，所以加索引能提高效率
                联合索引时，要注意sql 执行顺序，因为顺序可以导致某些联合索引失效（从B+树结构存放非聚簇索引的数据结构考虑原因）
            4.update delete 的where 条件列
            5.distinct 字段 需要创建索引
            6.多表join连接操作时，创建索引注意事项：
                连接表数量尽量不超过三张，
                对where条件创建索引
                对于连接的字段创建索引，并且该字段在多张表中类型必须一致，比如 要是int 都是int 
            7.使用列的类型小的值作为索引，原因是B+树的存储结构，降低树的高度，同时节省空间
            8.使用字符串前缀创建索引
                1，节省空间 2，节省比较时间   这种索引又称为前缀索引
            限制索引个数 一张表不超过6个
        不适合加索引的字段：
            1. 数据量小的数据表不适合加索引
            2. 大量重复数据的列上不适合创建索引
            3. 避免对经常更新的列增加索引
            4. 不建议用无序的值作为索引，原因是B+树的存储结构，在索引比较时需要转化为ascii,容易造成页分裂
            5. 删除很少用或者不用的索引
            6.不要定义冗余或重复的索引
                示例： 假设一张表中已经定义了联合索引(a,b,c),由于最左前缀原则（详细参看B+树的数据结构），这种索引对 查询条件 a , a,b, abc 同样有用，这时候就不需要单独创建a 或者 ab 的索引了
                设置联合索引时， 尽量将查询几率大的条件排前面，重复索引还需要注意，是否又重新定义了主键的索引，这种算是重复索引
    性能分析工具的使用：
        数据库服务器优化步骤：
            查看系统性能参数
                show (global|session) status like '参数'
                参数列表：
                    Connection: 连接Mysql服务器的次数
                    uptime: 服务器运行时常
                    Show_queries: 慢查询的次数（默认为10S,可以修改）
                    Com_select: 查询操作的次数
                    Com_insert: 插入操作的次数
                    Com_update: 更新操作的次数
                    Com_delete: 删除操作的次数
            统计sql 查询成本  last_query_cost
                show status like 'last_query_cost;'  返回值为查询数据所需要的数据页个数（B+树的数据结构）
                sql 查询是一个动态过程，从页加载的角度来看，
                    1.位置决定效率。如果查询页就在数据库的缓冲池中，效率最高。因为可以避免内存磁盘数据传输造成的时间开销
                    2.批量决定效率。如果对磁盘中单一页进行随机读，效率很低。如果采用顺序读取方式，批量对页进行读取，平均（注意这里是平均）一页的读取效率会提升很高
                所以读取数据时，首先考虑数据的存储位置，经常使用的就放到缓存池中，有以下几种策略：
                    LRU 策略是一种基于使用频率的策略。当缓冲池满时，最久未使用的数据页会被置换出去，新的数据页会进入缓冲池。这是 MySQL 默认的缓冲池管理策略。
                    LFU 策略是基于数据页使用频率的策略。它优先保留使用频率高的数据页，从而提高缓冲池中的数据命中率。
            定位查询 慢查询日志， 默认不开启
                查看当前是否开启
                    show variables like 'slow_query_log'    查看
                    set global slow_query_log = on;         开启
                    show variables like 'slow_query_log_file'    查看日志文件路径
                设置超时时间
                    show variables like '%long_query_time%'  查看当前阈值
                    set global long_query_time = 1;                设置阈值，注意这个变量是系统变量，同时也是会话变量 session
                    以上可以通过更改mysql 配置文件  systemctl restart mysql;
                查看慢查询记录：
                    show status like 'slow_queries' 返回慢查询的次数
                查找慢查询的sql语句
                    工具：
                        mysqldumpslow 
                        参数
                                -s ORDER     what to sort by (al, at, ar, c, l, r, t), 'at' is default
                                                al: average lock time
                                                ar: average rows sent
                                                at: average query time
                                                c: count
                                                l: lock time
                                                r: rows sent
                                                t: query time
                                -r           reverse the sort order (largest last instead of first)
                                -t NUM       just show the top n queries
                                -a           don't abstract all numbers to N and strings to 'S' 不加，查询条件int 型显示N
                                -n NUM       abstract numbers with at least n digits within names
                                -g PATTERN   grep: only consider stmts that include this string
                                -h HOSTNAME  hostname of db server for *-slow.log filename (can be wildcard),
                                            default is '*', i.e. match all
                                -i NAME      name of server instance (if using mysql.server startup script)
                                -l           don't subtract lock time from total time
                        查询到指定sql语句后，可以使用explain 进行分析
                关闭慢查询日志： 生产环境下需要关闭
                    set global slow_query_log = off;
                删除日志文件
                    调优结束后需要删除日志文件
                查看sql执行成本：cpu，IO等信息，根据参数可以设置
                    开启功能：
                        show variables like 'profiling';
                        set profiling=on;
                    show profiles; 显示所有sql 操作
                    show profile (id); 显示最新一条的查询数据
    explain 的使用（重点..........................................................）
        作用：
            可以帮助查看某个查询语句的具体执行计划，这里的计划是mysql 的优化器做的、
        输出参数:
            表的读取顺序
            数据读取操作的操作类型
            哪些索引可以使用
            哪些索引不可用
            表之间引用
            每张表有多少行被优化器查询
        返回字段含义
            id: 在一个大查询语句中每个select 都有一个ID
                id 相同，可以认为是一组，从上往下执行
                在所有组中，id值越大，优先级越高，越先执行
                id号每个号码，表示一趟独立的查询，一个sql查询趟数越少越好
            select_type: select 关键字对应查询的类型
                一条大的查询语句中可以包含若干个select关键字，每个select关键字代表一个小查询，每个select关键字可以包含若干张表。每一张表都对应执行计划输出的一条记录。根据select_type 属性，来推测扮演的角色

            table: 表名
                查询时用到几张表，就返回几条信息，包括可能存在的临时表
            partitions: 匹配分区信息
            type: 针对单表的访问方法(ALL)（重点）
                参数解释：
                    （从最好到最坏），优化的目标：最少到range级别，要求ref级别，最好是const级别
                    system: （最理想）
                        表中只有一条记录，存储引擎中统计数据是精确的，比如myisam memory,innodb 不支持
                    const:
                        根据主键或者唯一的二级索引列与常数(注意这个要求)进行等值匹配，对单表访问方式就是const
                    eq_ref:
                        当被驱动表通过主键或者唯一的二级索引列等值匹配的方式进行访问。主键或者唯一二级索引是联合索引
                    ref:
                        普通的二级索引与常量进行等值比较获取某个值
                    ref_or_null
                        当对普通二级索引进行等值匹配，该索引列可能是NULL时
                    unique_subquery:
                        针对一些包含IN子查询的查询语句中，如果优化器决定将in 子查询转换为exists子查询，并且子查询可以使用主键进行等值匹配的话
                    range:
                        使用索引获取某些范围区间的记录
                    index:
                        可以使用索引覆盖，但是需要扫描全部记录。
                    all:
                        全表扫描
            possible_keys: 可能用到的索引
            key: 实际用到的索引
            key_len: 实际用到的索引长度（联合索引）
                值越大越好
            ref: 当使用索引列等值查询时，与索引列进行等值匹配的对象信息
            rows: 预估的需要读取记录条数
                值越小越好，越小越有可能在同一页
            filtered: 经过过滤后剩下条数的百分比，
            Extra: 额外信息（重点）
                能够更准确的理解Mysql到底是如何执行给定的查询语句的
                no table use
                    语句中没有选中表
                impossible where
                    where查询条件永远为false
                using where
                    查询条件使用where
                using index
                    索引覆盖并且不用回表（推荐）
                using index condition（如果查询过程中要使用索引条件下推，则显示此消息，表示可以做优化）
                    搜索条件出现索引列，但是不能使用索引列
                    eg:
                        select * from T1 where k1 > 'z' and k1 like '%a'   其中k1 列上有索引
                        执行顺序，首先在k1 索引B+树上找到所有符合 k1 > 'z'的数据
                        然后回表找到所有数据的原始数据，再在结果中筛选满足 k1 满足 k1 like '%a'的数据
                using temporarry
                    说明当前语句执行时会创建临时表，应该尽量规避，比如指定索引来代替使用的临时表
    explain 
        不考虑各种cache
        不能显示mysql执行查询时所作的优化工作
        不会有关触发器，储存过程信息
        部分统计信息是估算的
    四种输出格式：
        传统模式
        json模式
            输出信息最详细
                explain format=json 
        tree格式
        可视化输出
    分析优化器执行计划， trace
        追踪优化器做出的各种决策
        功能默认关闭，
            开启并设置格式为json
            set optimizer_trace="enable=on",end_markers_in_json=on;
            set optimizer_trace_max_mem_size=10000  设置内存，防止内容过多显示不全
            开启后执行信息都存放到一张表中
            select * from information_schema.optimizer_trace \G
    mysql 监控分析视图  sys schema
        查询冗余索引
            select * from sys.schema_redundant_indexes;
        查询未使用的索引
            select * from sys.schema_unused_indexes;
        查询索引使用情况
            select index_name from sys.schema_index_statistics where table_schema='table_name'
        查询表的访问量
        占用bufferpool比较多的表
        监控sql执行频率等
索引优化与查询优化
    sql查询优化大致分为两个方面：
        物理查询优化
            索引和表连接方式
        逻辑查询优化
            sql的等价变化
    索引失效
        用不用索引，最终是优化器决定的，是基于开销而不是规则。sql 语句是否使用索引，和数据库版本，数据量，数据选择度都有关系。
        最佳左前缀规则
            注意where 查询顺序 和 生成联合索引的顺序，最佳左前缀
            可以为多个字段创建索引，一个索引最多包含16个字段，过滤条件必须按照索引建立的顺序，依次满足，一旦跳过某个字段，索引后面的字段无法使用。原理还是B+树存储索引的数据结构
        主键插入顺序
            建议是主键依次递增
        计算，函数，类型转换（特别要注意默认转换的情况和手动）导致索引失效
            特别小心字符串操作函数和隐式转换int
            where 条件存在表达式 fields +-*/ 等操作也会影响
        范围条件右边列表索引列失效
            where col1=1 and col2 > 10 and col3 =90 其中索引index(col1,col2,col3)中只用col1 和 col2
            其中where 后面三个字段的顺序，优化器会自动调优------------>，但是不能改变的索引定义的顺序
            优化方式：
                将查询条件涉及范围概率大的字段放到最后面，这样可以尽最大能力使用索引
            注意这里的右指的是索引定义的顺序
        不等于（!= or <>）索引失效
            理解的话，还是考虑B+树的结构
        is null 可以使用索引， is not null 无法使用索引
            类比于上方 != 
            一般设计数据表时，字段设置为not null 约束，比如int类型，设置默认值为0，字符串默认值为'',这样可以更高效利用索引
        like 以通配符% 开头的索引失效
            理解结合B+树来理解，上来不确定，在索引中无法确定，这里可能与索引中字符串比较算法有关
        阿里巴巴 规范：
            严禁在搜索中使用左模糊或者全模糊，
        or 前后存在非索引的列，索引失效
            理解结合B+树，由于or条件的特性，假如前后存在非索引字段，优化器将不能找到合适的索引来应用到该查询语句中，所以失效
        数据库和表统一使用utf8mb4 编码，这种编码方式兼容性最好。可以避免某些字符集转换产生乱码，不同字符集进行比较前需要转换会造成索引失效
    关联查询的优化
        内连接 vs 外连接（左外，右外等）
        以左外连接为例：
            连接条件两边数据类型要保持一致，否则会产生索引失效
        内连接
            两种连接方式有少许差别，其中外连接中，有明确的驱动表和被驱动表（但也有特殊情况），而内连接左右两张表角色相同，优化器会考虑最优情况进行索引设置，如果只有一张表有索引，则这张表会作为被驱动的表来处理
            在两个表连接条件都存在索引的情况下，小表作为驱动表
    join 底层原理：
        驱动表 和 被驱动表的概念
            优化器会转换
        sql 语句的执行顺序

        BNLJ 算法
    子查询优化：
        子查询可一次性完成很多逻辑上需要多个步骤才能完成sql操作
        子查询效率比较低的原因：
            执行子查询时，内层查询语句结果会生成一个临时表，消耗资源
            临时表中不存在索引，性能影响
            子查询返回的结果集，查询性能影响会有点大
    优化方式：
        可以使用join查询来替代子查询，连接查询不需要创建临时表，速度相对较快
        尽量不要使用not in 或者 not exists, 推荐方式用left join where is null 方式更高
    排序索引：
        在where 和 order by 的字段中都加上索引
        mysql 支持两种排序
            index
                索引可以保证数据有序性
            filesort
                一般存在内存中，占CPU较多
            注意点：
                order by 时不limit， 索引失效，有特殊情况，索引覆盖的情况，就是说，索引两个字段正好是查询的两个字段，不需要回表操作
        目的是where 子句中避免全表扫描，在order by 避免使用filesort 排序。
        order by 顺序错误，导致索引失效   最左前缀原则
        order by 规则不一致时，索引失效。主要有顺序错（单个字段，索引默认为升序），不索引，方向反（指的是多个字段的方向），不索引
        order by 无过滤，不索引 联合索引中，最左侧如果为等值条件，则索引会忽略后面字段
            filefort 效率不一定慢 index
                原因：
                    所有的排序都是在过滤之后才执行的（参考sql 语句的执行顺序），如果能过滤掉很大部分数据，剩下的排序不是特别消耗性能。即使添加了index,实际提升有限。
                    索引发挥最大性能在于遍历定位大量数据时，能引发最大性能。
                两个索引同时存在，优化器会选择最优化的方案
    gruop by 优化  
        大体和order by 相同， 即使过滤条件中没有用到索引，group by 仍可使用索引
        group by 先排序后分组，遵循最左前缀原则
        当无法使用索引列时，增大max_length_for_sort_data sout_buffer_size 参数设置
        where 效率高于having
        减少order by
        包含order by,group by,distinct 这些语句，where 过滤结果最好保持在1000行内

    limit 优化  例如这种  limit 2000000,10
    思路1：
        在索引上完成分页，根据主键然后查询
        select * from student s ,(select id from student order by id limit 200000,10) a where s.id = a.id
    思路2：
        select * from student where id > 200000 limit 10;   这里where 能过滤掉大部分数据
    优先考虑索引覆盖
        索引覆盖的意思是，查询字段正好时联合索引中的字段，避免回表，效率比较高
        上述规则部分会更改，例如存在<>的情况下索引失效，但是查询字段包含索引字段时，可以用
        % 在前时，同上，也可使用
        优势：
            避免回表
            可以把随机IO变为顺序IO提高效率
    索引下推（ICP）：
        场景：
            联合索引中，假设条件1在where的第一个，然后正常情况下执行顺序为先执行where第一个条件，然后执行回表，再然后执行后面的数据，当然这种情况一般出现在第二个条件为"%str"这种情况，导致索引失效。这时候假如说第一个条件过滤后条数过多，对应回表条数过多，效率低下。索引下推的作用就是将索引下推，先过滤，然后执行下一个过滤条件，最后回表
            eg: select * from T where age > 10 and name like '%zhang' and grand like '%abd%'
            index： create index index_name on T(age,name,grand)
        ICP 开启和关闭
            set optimizer_switch = 'index_condition_pushdown=on/off'  打开/关闭
        使用索引下推时，使用explain 进行性能分析时，extra 列内容显示的是user index condition

    其他查询优化策略
        exists in  区分
             基本原则是小表驱动大表
        count(*) 和 count(字段的效率)
            和数据库和使用的引擎有关系
            以mysql innodb 为例：
                count(*)  count(1) 效率相同
                count(field) count(1) 来比较，统计时间在底层实现都需要使用索引，假如存在二级索引，并且count(fields)为索引字段，并且二级索引比较小，可以认为两者用的同一个索引，效率一样
        关于select(*)
            会导致覆盖索引失效
            在sql语句预处理时也有部分性能损失，如果加上唯一索引，则没有性能提升
        limit 1
            针对全表扫描语句，如果确定结果只有一条，加上limit 1 会加快速度
        多使用commit
            释放资源
            锁
    淘宝数据库主键怎么设计
        主键设计应该考虑数据库层面和业务层面
            自增ID的问题
                安全性不高
                性能差，需要数据库服务端生成
                交互多，业务层还需要执行一次last_insert_id()函数才能知道刚刚的插入值
                局部唯一性。在表中唯一，在其他地方不唯一，不适合分布式
            根据业务字段作为ID
    sql 优化原则：
    mysql 范式
        数据库表的基本原则，规则称为范式
    数据表设计原则
        三少一多
        表个数少
        表字段个数少
        联合主键字段个数越少越好
        使用主键和外键越多越好
    数据库对象编写建议
    数据库调优措施：
        读写分离：
            主从模式
            双主双从模式
        拆分表
            冷热数据分离
            增加中间表
            优化数据类型
                优先选择符合条件最小类型
            水平分表
using 关键子，可以代替外连接中两张表同名字段相等的情况
脏写
脏读
不可重复读
幻读

四种隔离级别
	查看：（显示隔离级别）
		show variables like 'transaction_isolation'；
	设置：
		set global | session transcation_isolation=''
	read uncommit
		读未提交
		主要解决
	read commit
	repeatable read
事务日志：
	redo log 重做日志
		提交后的日志
	undo log 回滚日志
		回滚日志