# API Endpoints

/ (get) - index route, reads all tasks from database and returns a json response containing all tasks:
```
{
  message: "Welcome to Task Manager",
  tasks: [
    {
      description: "hello"
      due_date_time: "0002-02-01T02:02"
​      id: 1
​      status: "a"
​​​       title: "a"
    },
    ...
  }
}
```

/create (post) - submit tasks to the database using formdata where the accepted form values are:
"title"
"description" (optional)
"status"
"due_date_time"

/update (put) - takes url param of id and status as follows to update a task with the supplied id with the new status /update?id=X&status=newStatus

/delete (delete) - takes url param of id to delete associated task from database, /delete?id=3


# Run the application locally
## Backend 
1. ```cd backend```
2. ```go mod download```
3. ```go build -o backend```
4. ```./backend```

## Frontend
1. ```cd frontend```
2. ```touch .env && echo "PUBLIC_API_URL = \\"https://backend-production-93b4.up.railway.app\"" >> .env```
3. ```npm i```
4. ```npm run dev```
5. go to http://localhost:5173


# Live Preview
https://frontend-production-f438.up.railway.app/

# Dear Hiring Team
Thanks for viewing my submission, the assessment was fun!
