module github.com/yuu/lightbulb

go 1.13

require (
	github.com/gin-gonic/gin v1.6.3
	google.golang.org/grpc v1.31.1
	lightbulb.org/bto v0.0.0-00010101000000-000000000000
	lightbulb.org/defaults v0.0.0-00010101000000-000000000000
)

replace lightbulb.org/bto => ./bto

replace lightbulb.org/defaults => ./defaults
