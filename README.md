# electric-car-sharing
CNAD Assignment1

go get -u github.com/go-sql-driver/mysql
go get -u github.com/gorilla/mux
go get -u github.com/golang-jwt/jwt

How to test on postman
3.1 USER
1.REGISTER MEMBER
Select POST curl http://localhost:5000/api/v1/users/register
Under Header
Key: Content-Type
Value: application/json

click raw
under Body:
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone_number": "1234567890",
  "password": "securepassword123",
  "membership_tier": "Basic"
}

2.LOGIN
Select POST curl http://localhost:5000/api/v1/users/register
Under Body:
{
    "email": (created during register)eg. "john.doe@example.com"
    "password": (created during register)eg. "securepassword123"
}


3.TOKEN 
go run test/test_token.go
will display:
Generated Token for Testing: (token)

curl http://localhost:5000/api/v1/protected
under Headers
Key: Authorization
Value: Bearer (token)

----------------------------------------------------------------------------
3.2 VEHICLE
1.Car availability
curl http://localhost:5000/api/v1/vehicles


2.Booking
curl http://localhost:5000/api/v1/vehicles/book

under header:
Content-Type: application/json
under Body:
{
  "user_id": 1,
  "vehicle_id": 1,
  "start_time": "2024-12-11T10:00:00Z",
  "end_time": "2024-12-11T12:00:00Z"
}

3.Cancel booking
check sql and do
select * from reservations
check which registrationID is under Status Active, only Active can be cancelled
curl http://localhost:5000/api/v1/vehicles/cancel

under Body:
{
  "reservation_id": (if Active)eg. 2
}

------------------------------------------------------------------------------
3.3 Billing
1.Tier-based pricing
curl http://localhost:5000/api/v1/billing/calculate

under Body:
{
  "user_id": 1,
  "start_time": "2024-12-15T08:00:00Z",
  "end_time": "2024-12-15T12:00:00Z"
}

2.Real time
curl http://localhost:5000/api/v1/billing/realtime

under Body:
{
  "user_id": 1,
  "reservation_start": "2024-12-10T08:00:00Z"
}

3.Invoicing
curl http://localhost:5000/api/v1/billing/generate-invoice

{
  "reservation_id": 2,
  "user_id": 1,
  "amount": 50.00
}

