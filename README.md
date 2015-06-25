RabbitMQ Stress Tester
======================
A fork of [backstop/rabbit-mq-stress-tester](https://github.com/backstop/rabbit-mq-stress-tester) by [Backstop Solutions Group](https://github.com/backstop)

Building
--------

    mkdir /your/full/golang/project/root  # e.g., $HOME/goprojects
    export GOPATH=/your/full/golang/project/root
    go get -d github.com/martymacgyver/rabbit-mq-stress-tester  # Download-only
    go build  github.com/martymacgyver/rabbit-mq-stress-tester  # Build it!

Running
-------

    ./rabbit-mq-stress-tester [arguments]

      --server, -s "localhost"     hostname for RabbitMQ server
      --producer, -p "0"           number of messages to produce, -1 to produce forever
      --wait, -w "0"               number of nanoseconds to wait between publish events
      --consumer, -c "-1"          number of messages to consume. 0 consumes forever
      --bytes, -b "0"              number of extra bytes to add to the message payload (~50000 max)
      --concurrency, -n "50"       number of reader/writer goroutines
      --quiet, -q                  print only errors to stdout
      --wait-for-ack, -a           wait for an ack or nack after enqueueing a message
      --help, -h                   show help

Examples
--------

Consume messages forever using 50 concurrent goroutines:

    ./rabbit-mq-stress-tester --server localhost --consumer 0 --concurrency 50

Produce 1,000,000 messages of 10KB each, using 50 concurrent goroutines, waiting 100 nanoseconds between each message. Only print to stdout if there is a nack or when you finish:

    ./rabbit-mq-stress-tester --server localhost --producer 1000000 --concurrency 50 \
                              --bytes 10000 --wait 100 --quiet

Notes
-----

Interesting observations using localhost on Windows (may apply to Linux as well):
  - Using --quiet on the consumer (newly added option) as well as the producer improves throughput if the output device starts blocking (e.g., Windows command line). Monitored throughput via the RabbutMQ dashboard.
  - If 1 <= wait <= 100, the throughput was around 1000 msgs/s (used bytes == 1 -- setting bytes <= 10000 had no influence on this)
  - If wait == 0, the throughput was an order of magnitude higher (I saw around 13.3K msgs/s). CPUs (4 physical cores plus hyperthreading enabled) becomes ~90% saturated.
  - Increasing producer concurrency above 50 didn't improve things noticeably (it seemed to slightly decrease performance). In fact, anything above 1 had no effect. Possible bug in the producer code? The consumer definitely benefits from having higher concurrency) 
  - By running two producers (concurrency == 1) and one consumer (concurrency >= 3 due to overhead), I was able to sustain around 27.5K msgs/s, which saturated the CPU. That was about as fast as I could get this setup to go using default RabbitMQ settings.

Final test setup:

    Term1: rabbit-mq-stress-tester --server localhost --consumer 0 --concurrency 3 --quiet
    Term2: rabbit-mq-stress-tester --server localhost --producer 1000000 --concurrency 1 --bytes 1 --wait 0 --quiet
    Term3: rabbit-mq-stress-tester --server localhost --producer 1000000 --concurrency 1 --bytes 1 --wait 0 --quiet
    