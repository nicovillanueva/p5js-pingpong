# p5js-pingpong

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
    GET /
        - matchid
        --
        200 (text/html):
        - [match render]
    PUT /serve/public
        --
        200:
        - matchid
    PUT /serve/private
        - userid/jwt token
        - opponent id
        --
        200:
        - matchid
    PATCH /return
        - jwt auth token
        - matchid
        - code
        --
        200:
```