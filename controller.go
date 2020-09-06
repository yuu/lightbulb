package main

import (
        "os"
        "log"
        "strconv"
        "net/http"
        "time"

	"github.com/gin-gonic/gin"

	"lightbulb.org/bto"
        "lightbulb.org/defaults"
)

const (
        MODE_OFF = 0
        MODE_FULL_LIGHT = 1
        MODE_NIGHT_LIGHT = 2
)

type Status struct {
        handler defaults.Defaults

        Power int
        BrightnessFull int64
        BrightnessNight int64
}

func (s *Status) save() {
        s.handler.Save(s)
}

func (s *Status) load() {
        s.handler.Load(s)
}

type Controller struct {
	irClient bto.LightbulbController
        state Status
}

type intParam struct {
        ID int `uri:"id"`
}
type floatParam struct {
        ID float64 `uri:"id"`
}

func NewController(irCtrl bto.LightbulbController) *Controller {
        home := os.Getenv("HOME")
        var state Status
        handler := defaults.New(home + "/.config/io.flutia.lightbulb.toml")
        state.handler = handler
        state.load()

        return &Controller{irCtrl, state}
}

func (c *Controller) On(ctx *gin.Context) {
        if c.state.Power != MODE_FULL_LIGHT {
                c.irClient.On()
        }

        c.state.Power = MODE_FULL_LIGHT
        c.state.save()
        var ret int
        if c.state.Power == MODE_OFF {
                ret = 0
        } else {
                ret = 1
        }
        ctx.JSON(http.StatusOK, ret)
}

func (c *Controller) Off(ctx *gin.Context) {
        if c.state.Power != MODE_OFF {
                c.irClient.Off()
        }

        c.state.Power = MODE_OFF
        c.state.save()
        ctx.JSON(http.StatusOK, c.state.Power)
}

func (c *Controller) Status(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, c.state.Power)
}

func (c *Controller) SetBrightness(ctx *gin.Context) {
        const coe = 20 // 5段階

        bt := ctx.Query("brightness")
        b, err := strconv.ParseInt(bt, 10, 64)
        if err != nil {
                ctx.JSON(http.StatusBadRequest, gin.H{"msg": err})
                return
        }
        log.Printf("newBrightness: %v", b)

        current := c.state.BrightnessFull / coe
        future := b / coe
        log.Printf("current: %v, future: %v", current, future)
        if current < future {
                for i := current; i < future; i++ {
                        log.Printf("up...")
                        c.irClient.Up()
                        time.Sleep(100 * time.Millisecond)
                }
        }

        if current > future {
                for i := future; i < current; i++ {
                        log.Printf("down...")
                        c.irClient.Down()
                        time.Sleep(100 * time.Millisecond)
                }
        }

        c.state.BrightnessFull = b
        c.state.save()
        ctx.JSON(http.StatusOK, c.state.BrightnessFull)
}

func (c *Controller) GetBrightness(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, c.state.BrightnessFull)
}
