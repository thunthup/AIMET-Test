
# AIMET Calendar Technical Test

An API service for a calendar app 


## Run Locally

Clone the project

```bash
  git clone https://github.com/thunthup/AIMET-Test.git
```

Go to the project directory

```bash
  cd AIMET-Test
```

Install dependencies

```bash
  go mod download
```

Start PostgreSQL database

```bash
  docker compose up -d
```


Start the server

```bash
  go run main.go
```

generate random data (mockData.sql contain 2000+ random data)

```bash
  python3 mock.py
```

## Deployment

To deploy this project run

```bash
  go build
  ./aimet-test
```


## Running Tests

To run tests, run the following command

```bash
  go test -v -cover ./...
```


## Architecture

![App Screenshot](https://github.com/thunthup/AIMET-Test/blob/main/Architecture.png?raw=true)

## Database schema

![App Screenshot](https://github.com/thunthup/AIMET-Test/blob/main/Event%20Schema.png?raw=true)
## API Reference

#### Get events with filters

```http
  GET /api/events
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `start_date` | `date(YYYY-MM-DD)` | **Optional**. filter event that starts from the given date |
| `end_date` | `date(YYYY-MM-DD)` | **Optional**. filter event ending before and on the given date |
| `year` | `year (YYYY)` | **Optional**. filter event that happen in the given year (will overide start_date and end_date) |
| `month` | `month (MM)` | **Optional**. filter event that happen in the given month, year must also be given else month is ignored (will overide start_date and end_date) |
| `keyword` | `string` | **Optional**. filter event that contain the keyword (case sensitive) |
| `sort_order` | `string` | **Optional**. the events are sorted by date and time. sort order can either be "asc" or "desc". default is "asc"|


#### Get event

```http
  GET /api/events/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. ID of event to fetch |

#### Create event
```http
  POST /api/events/
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `title`      | `string` | **Required**. Title of the event   |
| `event_date` | `date(YYYY-MM-DD)` | **Required**. Date of the event|
| `start_time` | `time(01:35:00+07)` | **Required**. Start time of the event|
| `end_time` | `time(01:35:00+07)` | **Required**. End time of the event|


#### Update event
```http
  PUT /api/events/${id}
```
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. ID of event to update |
| `title`      | `string` | **Required**. Title of the event   |
| `event_date` | `date(YYYY-MM-DD)` | **Required**. Date of the event|
| `start_time` | `time(01:35:00+07)` | **Required**. Start time of the event|
| `end_time` | `time(01:35:00+07)` | **Required**. End time of the event|

#### Delete event

```http
  Delete /api/events/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. ID of event to delete |