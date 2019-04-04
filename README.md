# p5js-pingpong
Matches:
    a. between two users [private]
        one starter
        all users require login
        users can try to join the match, authenticating with their user_id, but the creator must approve
    b. open to any number of users [open]
        one starter
        only starting registered user required
    c. solo [solo]
        only one id allowed

Issues
    - GET /matches/:id?all
    - Response: [sketch,userid]

Ideas
    challenges - no; simplemente se puede negar el acceso al sketch privado a terceros
    objetivos/lineas generales predefinidas
    solo round
        basado en previos publicos
    export gif
    thumbnails
