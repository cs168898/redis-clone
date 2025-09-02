# redis-clone
Building the Redis in-memory database to gain a deeper understanding of Redis technology

## How to start
1. In an terminal, type ```go run .```, this will start the net/http server to listen to requests and accept commands.
2. Go to a separate terminal and type 'redis-cli' to enter the command line interface for redis. You should see an ip address as the prompt.

## Commands
1. ```PING```                       -> Returns PONG, typically used to test connections.
2. ```SET   [key] [value]```        -> Returns OK, used to set key value pairs. Data type: map[string]string{}.
3. ```GET   [key]```                -> Returns Value, used to find the value for specific keys.
4. ```HSET  [hash] [key] [value]``` -> Returns OK, typically used to set nested hashmap values. Data type: map[string]map[string]string.
5. ```HGET  [hash] [key] ```        -> Returns Value, typically used to obtain values inside nested hashmaps.
6. ```HGETALL [hash]```             -> Returns Value, typically used to list ALL key and values inside a specific hash.

## main.go
This file starts the connection to listen on port 6379 and receives requests from users to execute commands based on handlers then write back to the user using the write method defined in the ```writer.go``` file.

