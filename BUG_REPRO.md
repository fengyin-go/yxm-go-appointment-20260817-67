# Bug 是什么

预约记录的底层存储容器没有初始化，新建预约时向 nil map 写入，导致服务 panic。

# 如何触发

先创建科室、医生、排班和患者，然后调用创建预约接口或 `BookAppointment`。

# 错误信息

进程出现 panic，日志包含 `assignment to entry in nil map`，调用方看到服务器内部错误或连接中断。
