package bto

import (
        "context"
        "log"
        "time"
)

type LightbulbController interface {
        On() (*WriteResponse, error)
        Off() (*WriteResponse, error)
        // Mode() (*WriteResponse, error)
        Up() (*WriteResponse, error)
        Down() (*WriteResponse, error)
}

type Config struct {
        Light struct {
                On  []uint32 `toml:"on"`
                Off []uint32 `toml:"off"`
                Up  []uint32 `toml:"up"`
                Down []uint32 `toml:"down"`
        } `toml:"lights"`
}

type iRLightbulb struct {
        client IRServiceClient
        conf Config
}

func NewLightbulbController(client IRServiceClient, conf Config) LightbulbController {
	return &iRLightbulb{client, conf}
}

func (ir *iRLightbulb) On() (*WriteResponse, error) {
        return ir.write(ir.conf.Light.On)
}

func (ir *iRLightbulb) Off() (*WriteResponse, error) {
        return ir.write(ir.conf.Light.Off)
}

func (ir *iRLightbulb) Up() (*WriteResponse, error) {
        return ir.write(ir.conf.Light.Up)
}

func (ir *iRLightbulb) Down() (*WriteResponse, error) {
        return ir.write(ir.conf.Light.Down)
}

func (ir * iRLightbulb) write(data []uint32) (*WriteResponse, error) {
        ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
        defer cancel()

        req := WriteRequest{
                Frequency: 38000,
                Data:      data,
        }
        res, err := ir.client.Write(ctx, &req)
        if err != nil {
                log.Fatalf("could not write: %v", err)
                return nil, err
        }

        return res, nil
}
