## Fetch Rewards Receipt Processor Challange 

A webservice that fulfils the API described below. 

## Point Rules
Rules are defined as follows: 
* One point for every alphanumeric character in the retailer name.
* 50 points if the total is a round dollar amount with no cents. 
* 25 points if the total is a multiple of `0.25`
* 5 points for every two items on the receipt 
* If the trimmed length of the item description is a multiple of `3`, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned. 
* 6 points if the day in the purchase date is odd
* 10 points if time of the purchase is after 2:00 pm and before 4:00 pm 

## Usage

1. Clone repo to your desired directory 
1. Navigate to `fetch` directory 
1. To start the `server`, run `go run main.go` 
1. To gracefully shutdown the server, press `ctrl + C` 


## API Docs

### Process Receipt Endpoint

```http
  POST /receipts/process
```
#### body example 
```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```
#### Successful Response 
```json
{
  "id": "72d49ca8-c7bc-4af2-8adf-8bd694eaee65"
}
```
 
### Get Points Endpoint

```http
  GET /receipts/{id}/points
```

#### Successful Response
```json
{
 "points": 109
}
```
