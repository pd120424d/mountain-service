# mountain-service


# Test commands (cmd, Windows):

### Create an employee
```
curl -X POST http://localhost:8082/employees -H "Content-Type: application/json" -d "{\"username\": \"test-user\", \"password\": \"securepassword\", \"first_name\": \"Bruce\", \"last_name\": \"Lee\", \"gender\": \"M\", \"phone\": \"123-456-7890\", \"email\": \"test-user@example.com\", \"profile_picture\": \"/path/to/picture.jpg\", \"profile_type\": \"Medic\"}"
```

## Delete an employee (soft-delete)
```
curl -X DELETE http://localhost:8082/employees/1
```

## Get all employees (which are not deleted)
```
curl http://localhost:8082/employees
```


# Generating swagger documentation

Navigate to folder with `main.go` and then execute:

```
swag init -g main.go --pdl 3
```
### Explanation:
  - `-g main.go`: Specifies the entry point for the Swagger generation. This is the file where the main function and some Swagger annotations are located.
  - `--pdl 3`: Sets the package depth level to 3. This tells swag to scan directories up to three levels deep from the starting point (main.go) for Swagger annotations


# Running services with docker compose

```azure
# from root directory
docker compose up --build
```

# Accessing employee service API

```
http://localhost:8082/swagger/index.html
```