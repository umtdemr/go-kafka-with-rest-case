# go-kafka-with-rest-case

## Case

> Congratulations Ümit Demir, you have accessed the task content!

> Here are the requirements for the task we expect:

> A REST API should be developed using a language of your choice (preferably Java / Go / PHP / Python) to handle GET, POST, PUT, and DELETE requests with four different endpoints. The POST and PUT methods should accept requests with an empty request body.

> For each request received by this API, a successful response should be returned randomly between 0-3 seconds, and just before the response is returned, a log entry should be written to a log file with the content "{method type},{request processing time in ms},{timestamp}". Example log: "GET,1000,1614679220".

> An asynchronous job should take the lines written to the log file in real-time and send them to Kafka in a specified format.

> A consumer should catch the log information sent to Kafka and write it to a suitable database (such as RDBMS or NoSQL).

> A dashboard screen should graphically show how long the API requests in the last hour took to complete, and the graph should be updated live.

> Different request types (like PUT, DELETE) should be color-coded accordingly.
