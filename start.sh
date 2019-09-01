#!/bin/bash
#编写完脚本给脚本执行权限  chmod u+x test.sh

pwd
echo "hello,world."
echo "第1个参数：$1"
echo "第1个参数：$2"
echo "第1个参数：$3"
echo "第1个参数：$4"
echo "第1个参数：$5"
echo "获取传递的参数个数: $#"
echo "给脚本传递的所有参数: $@"
echo "当前脚本进程ID: $$"
#	$#:获取传递的参数个数
#	$@:给脚本传递的所有参数
#	$?:脚本执行完成之后的状态，失败>0 or 成功=0
#	$$:脚本进程执行之后对应的进程ID

#取值
#取普通变量的值
#ca=123  v=$ca  v=${ca}
#取系统变量的值
#   path=$(pwd) path=`pwd` 将命令pwd获取的值赋值给path
#-------if-----------
if [ -d $1 ];then
  echo "$1 是一个目录!"                                                
elif [ -s $1 ];then
  echo "$1 是一个文件, 并文件不为空"
else
  echo "$1 不是目录, 有肯能不存在, 或文件大小为0"
fi
#-------for----------
list=`ls`
for var in $list;do
  echo "当前文件: $var"       #取值输出
  echo '当前文件: $var'       #原样输出                                
done
# 定义函数
is_directory()
{
    # 得到文件名, 通过参数得到文件名
    name=$1
    if [ -d $name ];then
        echo "$name 是一个目录!"
    else
        # 创建目录
        mkdir $name
        if [ 0 -ne $? ];then
            echo "目录创建失败..."
            exit
        fi  
        echo "目录创建成功!!!"                                                                                         
    fi  
}