### Running app

Please export the variables included in .env file by running the commands in the terminal:
export DB_USER=user
export DB_PASSWORD=password
export DB_PORT=3306
export DB_NAME=getground
export DB_HOST=127.0.0.1

after that run: "go run main.go"

### Book a table
allows you to add a table with the seating capacity

```
POST /tables
body: 
{
    "seating_capacity": int
}
response: 
{
    "seating_capacity": 10
    "id": 2,
}
```

### Add a guest reservation to the list

allows you to the guests at the specified table, if there is insufficient space, the an error should be thrown

```
POST /reservation_list/name
body: 
{
    "table_id": int,
    "accompanying_guests": int
}
response: 
{
    "name": "string"
}
```

### Get the guest list

```
GET /reservation_list
response: 
{
    "guests": [
        {
            "name": "string",
            "table_id": int,
            "accompanying_guests": int
        }, ...
    ]
}
```

### Guest Arrives

A guest may arrive with the guests that are not included in the guest list.
If the table is expected to have extra space, allow them to come. 
Otherwise, this method should throw an error.

```
PUT /reservations/name
body:
{
    "accompanying_guests": int
}
response:
{
    "name": "string"
}
```

### Guest Leaves

All their accompanying guests leave as well, when a guest leaves.

```
DELETE /guests/name
response code: 204
```

### Get arrived guests

```
GET /reservations
response: 
{
    "guests": [
        {
            "name": "string",
            "accompanying_guests": int,
            "time_arrived": "string"
        }
    ]
}
```

### Count number of empty seats

```
GET /available_seats
response:
{
    "seats_empty": int
}
