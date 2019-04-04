## [React] Frontend

**!?**

## [Golang] API
Content-Type: application/json; unless noted
```
/users
    PUT /
        - username
        - password
        --
        200:
        - userid
    POST /login
        - username
        - password
        --
        200:
        - jwt bearing token
    
/matches
    GET /:matchId
        - matchid
        --
        200 (text/html):
        - [match render]
    PUT /serve
        playerOneID
        firstSketch
        mode: open/private
        --
        200:
        - matchid
    PATCH /return/:matchId
        - jwt auth token
        - matchid
        - code
        --
        200:
```