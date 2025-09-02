# redis-clone
Building the Redis in-memory database to gain a deeper understanding of Redis technology


## main.go
This file starts the connection to listen on port 6379 and receives requests from users to execute commands based on handlers then write back to the user using the write method defined in the ```writer.go``` file.