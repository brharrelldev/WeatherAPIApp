# Weather API

## Description

This is just a simple app that takes a city and state, then it gives the current weather

##Folder structure

| Folder Name  | Description                                                     |
|--------------|-----------------------------------------------------------------|
| Project Root | This just contains main and some utility function to get config 
| Service      | Creating the base http service                                  |
|models        | the modeling of request and response into structs|
| constants    | mostly for errors and custom errors, but could include enums
|validation    | to validate incoming data, only city and states are validated in this example


## Running Code

To build code simply run:

```make build```

To run code simply type:

```make run```

In case you need to run it with a debugger in your IDE (like delve), run:

```make debug-run```

Inside the `Makefile` there is `API_KEY`.  Replace this with your actual API key before running

How to make a request:

```curl --request POST \
  --url http://localhost:3000/weather \
  --header 'Content-Type: application/json' \
  --data '{
	"city": "Atlanta",
	"state": "georgia"
}
```


Response will look like this:

```
{
	"feeling": "HOT!",
	"temperature": 81.66
}
```
