# Jopit API Items

This is the API for the Items Domain .

To run swagger follow this steps:

To build documentation run the following steps:
     
1. Ensure to have swagger installed: 

     ```go install github.com/swaggo/swag/cmd/swag@v1.8.12```

     If you find some difficulties with it try cleaning deps and re install:
        
     - ```go clean -i github.com/swaggo/swag/cmd/swag```
     - ```go clean -i github.com/swaggo/swag```
     - ```go clean -modcache```

2. Set up swagger env variables on your terminal:

     - ```export PATH=$PATH:$(go env GOPATH)/bin```

3. Copy the `handlers/` pkg to `/src/main/api/` temporarily:

     - ```cp -r ./src/main/domain/handlers/ ./src/main/api/```
    
4. Execute command to build the docs folder on the root of the project:

     - ```swag init --parseDependency --parseDepth 5 -d ./src/main/api/```

5. Delete the copied main.go and fix any failed import:

     - ```rm -rf ./src/main/api/handlers/```

6. Run on debug mode and access on localhost to ensure it worked:

     - ```http://localhost:8080/items/swagger/index.html```

## Testing

You can run tests separatelly by branches and completely:

Handler tests:      

- ```go test -v src/tests/internal/domain/handlers/items/items_test.go```
- ```go test -v src/tests/internal/domain/handlers/categories/categories_test.go```

Service tests:      

- ```go test -v src/tests/internal/domain/services/items/items_test.go```
- ```go test -v src/tests/internal/domain/services/categories/categories_test.go```

Repository tests:  

- ```go test -v src/tests/internal/domain/repositories/items/items_test.go```
- ```go test -v src/tests/internal/domain/repositories/categories/categories_test.go``` 

(on Ubuntu 22.04 env)

Test full project and generate coverage: 

1. Run the test for the whole project and generate the coverage.out file on the root of the project:

     - ```go test -v -coverprofile=./coverage.out -covermode=atomic -coverpkg ./src/main/domain/... ./...```

2. Build the html coverage file on the root of the project:

     - ```go tool cover -html=./coverage.out -o ./coverage.html```
     
3. Drop the coverage.html on any html online complier an check for the coverage. 

4. Delete the generated coverage files before pushing the changes: 

     - ```rm -f coverage.out && rm -f coverage.html```

## Telemetry

At the moment Traces and Logs has been enabled to track and push to Grafana Cloud.  

## MercadoLibre auth flow (plain English)

Here is the short, no-code version of how the MercadoLibre OAuth flow works in this project:

1. The user clicks "Connect MercadoLibre" in the UI.
2. The frontend asks the ETL API for the MercadoLibre login URL.
3. The ETL API builds that URL and returns it.
4. The frontend redirects the user to MercadoLibre to log in and approve access.
5. MercadoLibre redirects the user back with a temporary code.
6. The frontend sends that code to the ETL API.
7. The ETL API exchanges the code for an access token and refresh token.
8. The ETL API stores those tokens in MongoDB.
9. When we need to call MercadoLibre, we read the token from MongoDB.
10. If the token is close to expiring, the API refreshes it automatically and saves the new one.

In short: the browser never stores tokens, only the backend does, and it keeps them fresh for you.
