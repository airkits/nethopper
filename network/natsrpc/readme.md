## 连接nats集群

选择一台nats连接
a服务 a0 a1 a2
b服务 b0 b1 b2
同时连接到nats

stream创建6条
a0 a1 a2
b0 b1 b2
a0 sub ag  b[0-2].a0
a1 sub ag  b[0-2].a1
a2 sub ag  b[0-2].a2

b0 sub bg  a[0-2].b0
b1 sub bg  a[0-2].b1
b2 sub bg  a[0-2].b2

逻辑处理 a0<->b0
a0 pub a0.b0
b0 pub b0.a0

逻辑处理 a0广播->b0 b1 b2
a0 pub bg

逻辑处理 a0广播->a0 a1 a2
a0 pub ag

逻辑处理 b0广播->a0 a1 a2
b0 pub ag

逻辑处理 b0广播->b0 b1 b2
b0 pub bg


client <-10x2-> nginx <-5x2-> logic <-5x2-> nats <-5x2-> gamedb <-> mysql