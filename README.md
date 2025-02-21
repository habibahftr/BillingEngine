```bash
git clone https://github.com/habibahftr/BillingEngine.git
go mod tidy

```set up database
CREATE DATABASE billing_engine

```start server
go run main.go

```API endpoints
POST ...loan/create        --> create loan
GET ...loan/{CUSTOMER_ID}  --> get outstanding data
POST ...loan/payment       --> make payment
