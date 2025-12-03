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
