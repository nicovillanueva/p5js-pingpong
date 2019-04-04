module github.com/nicovillanueva/p5js-pingpong

go 1.12

require (
	github.com/labstack/echo/v4 v4.0.0
	github.com/labstack/gommon v0.2.8
	github.com/mitchellh/go-homedir v1.1.0
	github.com/nicovillanueva/p5js-pingpong/api v0.0.0-20190404130543-55234abf6c11
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
)

replace github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go v3.2.1-0.20180921172315-3af4c746e1c2+incompatible