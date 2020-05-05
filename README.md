# Drawy draw (name not final lol)
Needs some description

## Local development
### Client
- Enter the client directory
- Run `npm run dev`
- Client should be available at `localhost:8080`
### Server
- Enter the server directory
- Run `go run main.go`
- Service should be available at `localhost:3000`

## Production
Code is hosted at https://drawydraw.herokuapp.com

To push master to production:
```
git push heroku master
```

To push a different branch to production:
```
git push heroku $BRANCH:master
```
