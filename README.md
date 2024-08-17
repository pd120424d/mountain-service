# mountain-service


## Test commands (cmd, Windows):

### Create an employee

```
curl -X POST http://localhost:8082/employees -H "Content-Type: application/json" -d "{\"username\": \"jdoe\", \"password\": \"securepassword\", \"first_name\": \"John\", \"last_name\": \"Doe\", \"gender\": \"M\", \"phone\": \"123-456-7890\", \"email\": \"jdoe@example.com\", \"profile_picture\": \"/path/to/picture.jpg\", \"profile_type\": \"Medic\"}"
```