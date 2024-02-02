# RSS Aggregator

> Simple web server built with Go and PostgreSQL, that constantly checks to see new posts through RSS Feed and fetch them to the database, with concurrency.

This project was based in [FreeCodeCamp tutorial on Golang](https://youtu.be/un6ZyFkqFKo). 

### Functionality

The code creates a CRUD API where you can:
- Create a new User;
- Fetch users;
- Create new Feeds;
- Fetch feeds;
- Follow, as a user, a feed;
- Fetch the feeds you follow;
- Stop following a feed;
- Automatically, store new posts of that feeds in database.