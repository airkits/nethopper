rm -rf ~/natsdata/*
nohup nats-server -c nats1.conf &
nohup nats-server -c nats2.conf &
nohup nats-server -c nats3.conf &
