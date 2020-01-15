# creek_ext
本代码库用于介绍Creek使用案例，扩展Creek的自定义函数，经验技术交流等

Creek是一个轻量级的通用流式计算框架，作业内存消耗<10MB，运行时零依赖。Creek提供与Apache Flink完全兼容的SQL接口，采用申明式的作业定义，在云端静态编译，以可执行文件方式下发。

Creek在线作业编辑生成，请访问: http://creek.baidubce.com/


## In English
This Repo demos use case of Creek, as well as sharing UDF, and also serves as a place for communication.

Creek is a *super lightweight* general purpose stream processing framework, it's implemented in Go. It provides Apache *Flink* compatible SQL interface, user declares jobs in a JSON document, no coding is required.

Creek takes *less than 10MB* of memory, hence it can run on most of the resource constrained devices, e.g. edge devices.

Jobs are statically linked to executables, so it has *no dependency* during runtime, therefore deploying it is as simple as download and launch. 

To build a Creek job, refer to http://creek.baidubce.com/