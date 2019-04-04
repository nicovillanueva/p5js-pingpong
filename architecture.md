## Repo organization
Basically a copycat of perkeep: https://github.com/perkeep/perkeep

## [React] Frontend

**!?**

## CLI client

### Usage
${cmd} return ${matchid} ./path/to/sketch.js  # return a hit in a match
${cmd} serve ./path/to/sketch.js              # start a new match with a sketch
${cmd} join ${matchId}   # join a match

### Config
- user id

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