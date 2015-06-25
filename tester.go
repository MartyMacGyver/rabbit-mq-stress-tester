package main

import (
	"github.com/codegangsta/cli"
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
	"fmt"
)

var totalTime int64 = 0
var totalCount int64 = 0

type MqMessage struct {
	TimeNow        time.Time
	SequenceNumber int
	Payload        string
}

func main() {
	app := cli.NewApp()
	app.Usage = "RabbitMQ Stress Tester"
	app.Version = ""
	app.HideVersion = true
	app.HideHelp = true  // Hide help as a full command
	app.Author = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{"server, s", "localhost", "hostname for RabbitMQ server", ""},
		cli.IntFlag{"producer, p", 0, "number of messages to produce (-1 to produce forever)", ""},
		cli.IntFlag{"consumer, c", -1, "number of messages to consume (0 consumes forever)", ""},
		cli.IntFlag{"wait, w", 0, "number of nanoseconds to wait between publish events", ""},
		cli.IntFlag{"bytes, b", 0, "number of extra bytes to add to the message payload (~50000 max)", ""},
		cli.IntFlag{"concurrency, n", 50, "number of reader/writer goroutines", ""},
		cli.BoolFlag{"wait-for-ack, a", "wait for an ack or nack after enqueueing a message", ""},
		cli.BoolFlag{"quiet, q", "print only errors to stdout", ""},
		cli.BoolFlag{"help, h", "show help", ""}, // Retain --help/-h as switches
	}
	app.Action = func(c *cli.Context) {
		runApp(c)
	}
	app.Run(os.Args)
}

func runApp(c *cli.Context) {
	uri := "amqp://guest:guest@" + c.String("server") + ":5672"

	if c.Int("consumer") > -1 && c.Int("producer") != 0 {
		fmt.Println()
		fmt.Println("!! Error: Cannot specify both producer and consumer options together")
		fmt.Println()
		cli.ShowAppHelp(c)
		os.Exit(1)
	} else if c.Int("consumer") > -1 {
		fmt.Println("Running in consumer mode")
		config := ConsumerConfig{uri, c.Bool("quiet")}
		makeConsumers(config, c.Int("concurrency"), c.Int("consumer"))
	} else if c.Int("producer") != 0 {
		fmt.Println("Running in producer mode")
		config := ProducerConfig{uri, c.Bool("quiet"), c.Int("bytes"), c.Bool("wait-for-ack")}
		makeProducers(config, c.Int("concurrency"), c.Int("producer"), c.Int("wait"))
	}
}

func MakeQueue(c *amqp.Channel) amqp.Queue {
	q, err2 := c.QueueDeclare("stress-test-exchange", true, false, false, false, nil)
	if err2 != nil {
		panic(err2)
	}
	return q
}

func makeProducers(config ProducerConfig, concurrency int, toProduce int, wait int) {
	taskChan := make(chan int)

	for i := 0; i < concurrency; i++ {
		go Produce(config, taskChan)
	}

	start := time.Now()

	for i := 0; i < toProduce; i++ {
		taskChan <- i
		time.Sleep(time.Duration(int64(wait)))
	}

	time.Sleep(time.Duration(10000))

	close(taskChan)

	log.Printf("Finished: %s", time.Since(start))
}

func makeConsumers(config ConsumerConfig, concurrency int, toConsume int) {
	doneChan := make(chan bool)

	for i := 0; i < concurrency; i++ {
		go Consume(config, doneChan)
	}

	start := time.Now()

	if toConsume > 0 {
		for i := 0; i < toConsume; i++ {
			<-doneChan
			if i == 1 {
				start = time.Now()
			}
			log.Println("Consumed: ", i)
		}
	} else {
		for {
			<-doneChan
		}
	}

	log.Printf("Done consuming! %s", time.Since(start))
}
