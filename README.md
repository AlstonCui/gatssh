# gatssh

    gatssh是一个web ssh工具，基于go语言和beego框架和，支持同时向多台主机发送命令，并返回结果。类似ansible的shell模块。
    
    暂时只支持用户名与密码的形式连接主机；
    
    线程池暂时为1000，初始密码 admin:123456，可于setting页面修改。
    
    已经尝试在1000并发的情况下，操作15000台服务器，速度远优于ansible。
    
    支持执行结果的筛选和excel下载。

    
