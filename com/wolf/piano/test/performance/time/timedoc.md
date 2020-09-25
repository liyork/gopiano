场景：如何分析程序的运行时间与CPU利用率情况？
time指令
time go run test2.go 
3个指标
real：从程序开始到结束，实际度过的时间；
user：程序在用户态度过的时间；
sys：程序在内核态度过的时间。
一般情况下 real >= user + sys，因为系统还有其它进程(切换其他进程中间对于本进程回有空白期)。


/usr/bin/time指令
比time详细
/usr/bin/time -v go run test2.go  
还包括了：
CPU占用率；
内存使用情况；
Page Fault 情况；
进程切换情况；
文件系统IO；
Socket 使用情况；