## 连接nats集群

选择一台nats连接
a服务 a0 a1 a2
b服务 b0 b1 b2
同时连接到nats

stream创建6条
a0 a1 a2
b0 b1 b2
a0 sub a0.a0.* a[0-2].aa b[0-2].ba
a1 sub a1.a1.* a[0-2].aa b[0-2].ba
a2 sub a2.a2.* a[0-2].aa b[0-2].ba

b0 sub b0.b0.* a[0-2].b0.* b[0-2].bb a[0-2].ab
b1 sub b1.b1.* a[0-2].b1.* b[0-2].bb a[0-2].ab
b2 sub b2.b2.* a[0-2].b2.* b[0-2].bb a[0-2].ab

逻辑处理 a0<->b0
a0 pub a0.b0.*
b0 pub a0.a0.*

逻辑处理 a0广播->b0 b1 b2
a0 pub a0.ab

逻辑处理 a0广播->a0 a1 a2
a0 pub a0.aa

逻辑处理 b0广播->a0 a1 a2
b0 pub b0.ba

逻辑处理 b0广播->b0 b1 b2
b0 pub b0.bb


client <-10x2-> nginx <-5x2-> logic <-5x2-> nats <-5x2-> gamedb <-> mysql