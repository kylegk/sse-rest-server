# SSE REST Server

## Table of Contents

1. [Description](#description)
2. [Methods](#methods)
3. [Getting Started](#getting-started)

## Description

In this project, you will find a very simple REST server that provides exam data collected from a third-party [Server-Sent Events](https://www.w3.org/TR/2015/REC-eventsource-20150203/) service. This REST server collects SSE messages in real time, stores these events in an in-memory database and returns this data to consumers in a JSON format.

The project makes heavy use of two libraries to collect and store the events. For retrieving events from the SSE service, the project uses [r3labs/sse](https://github.com/r3labs/sse) made by R3 Labs, and for storing the event data the project relies on the [go-membdb](https://github.com/hashicorp/go-memdb) in-memory database solution created by HashiCorp.

## Methods

**All Students**

```
/students
```

> Method: **GET**

> Lists all students that have received at least one test score

```
{
   "students" : [
      "Abdul_Emard",
      "Alexys.Price",
      "Andreane1",
      "Brady62",
      "Brielle.Balistreri85",
      "Camron.Stark46",
      "Claire36",
      "Clemens.Fahey",
      "Gunnar_Ledner",
      "Imani.Pagac93",
      "Isobel_Gottlieb46",
      "Jovanny.Gibson",
      "Makenna.Jacobson14",
      "Mckayla84",
      "Nels.Wehner",
      "Nelson56",
      "Ricky_Lesch",
      "Vance_Powlowski63",
      "Wilma_Kulas",
      "Zack20"
   ]
}
```

**Student**

```
/students/{id}
```

> Method: **GET**

> Lists the test results for the specified student, and provides the student's average score across all exams

```
{
   "average" : 0.70000000000000
   "exams" : [
      {
         "exam" : 15849,
         "score" : 0.65000000000
      },
      {
         "exam" : 15850,
         "score" : 0.75000000000
      },
   ],
   "student" : "Zack20"
}
```

**All Exams**

```
/exams/all
```

> Method: **GET**

> Lists all the exams that have been recorded, along with the student and score

```
{
   "exams" : [
      {
         "exam" : 15872,
         "score" : 0.757167038802041,
         "studentid" : "Abdul_Emard"
      },
      {
         "exam" : 15872,
         "score" : 0.778255850930371,
         "studentid" : "Alexys.Price"
      },
      {
         "exam" : 15872,
         "score" : 0.780391071563619,
         "studentid" : "Andreane1"
      },
}
```

**Unique Exams**

```
/exams
```

> Method: **GET**

> Lists all the unique exams that have been recorded

```
{
   "exams" : [
      15872,
      15936,
      15873,
      15937,
      15874,
      15938,
   ]
}
```

**Exam**

```
/exams/{id}
```

> Method: **GET**

> Lists all the results for the specified exam, and provides the average score across all students

```
{
   "average" : 84.666666666666667,
   "exam" : 15872,
   "scores" : [
      {
         "score" : 0.750000000000,
         "student" : "Abdul_Emard"
      },
      {
         "score" : 0.800000000000,
         "student" : "Alexys.Price"
      },
      {
         "score" : 0.990000000000,
         "student" : "Andreane1"
      }
}
```

**Add Exam**

```
/exams
```

> Method: **PUT**

> Add an exam to the data store

> `Request:`

```
{
        "exam": 12345,
        "score": 0.78,
        "studentid": "test.student"}
}
```

> `Response:`

```
{
        "message":"Succesfully added exam: 12345"
}
```

**Delete Exams**

```
/exams/{id}
```

> Method: **DELETE**

> Deletes all exams in the datastore with a matching exam id

> `Response:`

```
{
        "message":"Successfully deleted {count} exams"}
}
```

## Getting Started

This project can either be built manually or run in a Docker container.

### Building manually

Building the project manually requires that you have a recent version of Golang installed on your system. You will also need to set the following environment variables in your shell prior to running the binary:

1. `SSE_SERVER_URL`: The url of the SSE server publishing event messages. You can use any url you wish, but the expectation is that events returned from the SSE server will be in the format: `{"exam": int, "studentId": string, "score": float}`

2. `APPLICATION_PORT`: The port number this server will listen on. You must include the "`:`" when assigning a port.

To build the project manually, perform the following steps:

```
export SSE_SERVER_URL="http://live-test-scores.herokuapp.com/scores"
export APPLICATION_PORT=":8080"
cd /path/to/project/
go build -o bin/goapp
cd bin/
./goapp
```

### Docker

This project includes a Dockerfile for easy building and containerization of the application. As with the manual process described above, you can use a different SSE server and port by modifying the environment variables in the Dockerfile.

To build the Docker container, perform the following steps:

```
cd /path/to/project
docker build -t launchdarkly-coding-test .
docker run -d -p 8080:8080 sse-rest-server
```

Once the container has started, you can connect to it and view the output of any logs in the application, such as requests and errors.
