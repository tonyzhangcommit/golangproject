常见命令：
    show create table tablename;
        导入数据
            source 文件路径， 命令行操作
        别名的使用(三种方式)：
            列名后直接写
            AS 关键字 + 别名
            列名 + 双引号包裹别名
        去重：
            distinct 关键字 DISTINCT + 字段名
            注意放置顺序， 使用为单独使用
            可以使用CONCAT() 关键字将多个字段组合进行多字段组合去重查询（效率低下）
        空值参与运算
            null 不等于0，''
            空值参与运算时，结果一定为空值null，参与运算时，注意使用内置函数转换下，IFNULL(fields,0)
        着重号``：
            为了防止表名（字段名）和mysql 关键字重复的情况，eg:select * from `order`;
        查询常数：
            查询时，增加某个列的常数，eg:select '公司名',name from table_name;
        显示表结构：
            describe 简写 desc
        算术运算符
            +(存在隐式转换 eg:100 + "1")
            -
            * 
            /(div)
            %(mod) 注意符号只和被取模数相同 12 % -5 ==> 2 -12 % 5 ==> -2
        比较运算符， 
            注意，存在不能隐式转换时，转换为0，只适用于两边一边为整型，一边为字符串的情况
            null 不等于 null
            =
            <==>  安全等于,可以对null 进行判断，两边都是null时，返回1，只有一边为null，返回0
            <>  !=
            <=
            >=
            >
            <
            is null  功能等同于<==>
            is not null  
            isnull() 一般为函数使用
            least()  最小， least(fields1,fields2...),注意不是面对列的
            greatest() 最大  同上
            betweent condition1 and condition2
            in (set)  eg: where id in (10,20,30)
            not in (set) 同上
            like 模糊查询  % 代表0或者多个不确定的字符， _ 代表一个不确定字符 \转义字符
            regexp  \ rlike
            not !
            and &&
            or ||
            xor  逻辑异或 两边不同为真
            asc  升序
            desc 降序
            order by fields  （不限于列名，还可以使用别名）
            （列的别名只能在order by 语句中，不能在where语句中）
            注意sql 语句的执行顺序
            eg: select * from users order by age desc/asc   默认为升序排序
            二级排序示例：
            select * from users order by age desc,score asc;
            limit 分页用法
            limit 0，20， 第一个参数为偏移量，第二个为偏移个数（每页条数）
            limit  (pagesize -1) * pagenumber ,pagenumber
            mysql 8.0 offset 后面加偏移量多少
            eg: limit 2 offset 20  取 21 22 数据
            多表查询：
                select name,city from user,tcity from users,city where users.cityid = tcity.id   92 语法，推荐使用99语法 join table_name on conditions
                可以给表加别名
            等值连接 vs 非等值连接
                等值与否 是相对连接条件而言
            自连接 vs 非自连接
                自连接 与非自连接 相对于连接的表是不是自己
                select userid,username,manageid,managename from users u,users m where ... 92 语法
                select userid,username,manageid,managename from users u join users m on ... 99 语法
                把一张表看成两张表来处理
            内连接 vs 外连接
                把两个表看成集合，内连接为两张表的交集
                外连接
                    左外  左表的所有，加上右表符合条件的记录  left join
                    右外  右表的所有，加上左表符合条件的记录  right join
                    满外    
                union 
                union all  都是求并集(相对于union, 不去重，性能较好，推荐使用)
                实现要求，列数量一样，每个列的字符类型也一样
                eg:
                    select userid,cityid from users u left join city c on u.cid = c.cityid
                    union all 
                    select userid,cityid from users u right join city c on u.cid = c.cityid where u.uid is null;
                在on 条件关联后仍然可以加条件，实现进一步筛选操作
                using 关键字，可以将on 关键字后面的同名（不同表）的字段合并为一个
                charlength()
                length()
                concat()
                concat_ws('分隔符','str1','str2'....)
                upper()
                lower()
                left(str,number)
                right(str,number)
                ltrim()
                rtrim()
                trim()
                trim(s1 from s) 去掉字符串s开始与结尾的s1
                curdata() 返回当前日期
                curtime() 返回当前时间
                now()返回当前日期加时间
                unix_timestamp()
                from_unixtime()  时间戳转为普通格式时间
                还有一些函数用于返回年月日，季度，星期几等
                查询表达式if 语句用法
                case when then when then end
                聚合函数(作用一组数据，并对一组数据返回一个值)
                    只适用于数值类型
                    avg()  select avg(age) from users;
                    sum()  select sum(age) from users;
                    可以用于字符串，比较的是ASCII的值
                    也可以用于日期时间类型
                    max()  select max(age) from users;
                    min()
                    计算指定字段查询的个数，eg:count(userid)
                    count(1) 和 count(0)
                    innodb 引擎中，count(*) count(1) 直接读行数，时间复杂度为o(n),但是好于count(列名)
                    count()

                    常结合 group by  和 having 结合来使用
                    group by 后面加的字段，select 后面不一定需要，但是select 后面非函数调用字段，group by 必须有
                    如果过滤条件中使用了聚合函数（avg,min,max,sum等）， 必须用having 进行过滤
                    建议在having 中只使用涉及聚合函数相关条件的过滤，原因是sql的执行顺序
                    sql 执行顺序：
                        from table_name join on ... where group by  having select distinct order by limit
                        每个过程都产生一个虚拟表， 然后将虚拟表传到下一个语句中，可以将此特性作为sql 优化的一个方面
                    约束定义分类以及实现方法
                    查看约束
                        information_schema.table_constraints
                    not null
                    unique  特例：允许
                    主键
                    外键
                    默认值
                    选择性约束
                    视图介绍：
                        可以理解为一张虚拟表，用于将复杂查询的结果生成一张表，方便后续调用
                        还有一个用途是，方便某些权限对外放开的使用，针对不同用户指定不同的查询视图
                    创建方式：
                        create view view_name (视图字段别名)  as select fields1,fields2,... from ...
                        使用视图更新数据时，注意视图中字段在基表中是否存在（因为视图中有些字段是通过聚合函数计算出来的，基表中并不存在）
                    存储过程：
                        一组经过编译的sql语句的封装，保存在数据库服务器，可以通过客户端进行调用， 可以增加参数，也可以有参数返回
                        优势：
                            create procedure 存储过程名 （IN OUT ） in 代表输入参数 out 代表返回值
                            begin
                                存储结构体
                            end
                            编写时注意语句结尾; 有时候需要转义

                            调用
                                call 存储过程()
                    游标：
                        暂时跳过
                    触发器：
                        适用于存在关联表时，一张表信息发生改变，另外一张表信息同时需要改变，这里使用触发器保证数据一致性
                        当然可以将两步操作封装成一个事务，也可以保证数据一致性
                        触发器是由事件触发，事件通常为 insert update delete 
                        create trigger 触发器名称
                        {before|after} {insert|update|delete} on 表名
                        for each row
                        触发器执行语句块 可以是单条， 也可以是 begin end 包含的语句块
                        更新的数据记录在语句块中可以使用new 关键字来表示
                        create trigger test
                        before insert on users
                        for each row
                        begin
                            declare vard int
                            select age into vard from users where name = new.name
                            if new.name = vard then 
                                ...
                                endif
                        end //  这里需要用delimiter 提前转义下
                    查看/删除 触发器
                    show trigger    查看当前数据库中的触发器
                    show create trigger 触发器名称  查看某一个触发器的定义
                    drop trigger 触发器名称  一般结合if not exists


                        



            






        

