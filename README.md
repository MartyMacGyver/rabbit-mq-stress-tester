RabbitMQ Stress Tester
======================

Building
--------

    mkdir /your/full/golang/project/root  # e.g., $HOME/goprojects
    export GOPATH=/your/full/golang/project/root
    go get -d github.com/martymacgyver/rabbit-mq-stress-tester  # Download-only
    go build  github.com/martymacgyver/rabbit-mq-stress-tester  # Build it!

Running
-------

    ./rabbit-mq-stress-tester  - Make the rabbit cry

    USAGE:
        rabbit-mq-stress-tester [global options] command [command options] [arguments...]

        --server, -s "localhost"     hostname for RabbitMQ server
        --producer, -p "0"           number of messages to produce, -1 to produce forever
        --wait, -w "0"               number of nanoseconds to wait between publish events
        --consumer, -c "-1"          number of messages to consume. 0 consumes forever
        --bytes, -b "0"              number of extra bytes to add to the RabbitMQ message payload. About 50K max
        --concurrency, -n "50"       number of reader/writer Goroutines
        --quiet, -q                  print only errors to stdout
        --wait-for-ack, -a           wait for an ack or nack after enqueueing a message
        --help, -h                   show help
        --version, -v                print the version

Examples
--------

Consume messages forever:

    ./tester --server localhost --consumer 0

Produce 100,000 messages of 10KB each, using 50 concurrent goroutines, waiting 100 nanoseconds between each message. Only print to stdout if there is a nack or when you finish.

    ./tester --server localhost --producer 100000 --bytes 10000 \
             --wait 100 --concurrency 50 --quiet
