# Scheduler

The scheduler allows to schedule a task at a specific date thanks to a http api.  
The scheduler is backed by sqlite3 database, but other database drivers will be added in the future. 

## Task

A task is defined by a date and a list of publishers.
The execution of a task will trigger the list of publishers at the provided date.

## Publisher

A Publisher is a go plugin which will be executed at the execution of a task. 

### Build a publisher plugin

go build -o plugins -buildmode=plugin publisher/http/http.go

## Config



