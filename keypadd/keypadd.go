package main

import (
    "flag"
    "fmt"
    "log"
    "time"

    "github.com/kylelemons/gousb/usb"
)

var (
    debug = flag.Int("debug", 3, "Turn debug on")
)

func find_keyboard(ctx usb.Context) []*usb.Device {

    devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
        return usb.Class(desc.Class) == usb.CLASS_HID;
        })

    if err != nil {
        log.Fatalf("list: %s", err)
    }

    return devs

}


func get_keys(keyboard_in *usb.Endpoint, b *[]byte) {
    kb := *keyboard_in

    _, err := kb.Read(*b)

    if err != nil {
        log.Fatalf("Exploded on reading keys")
    }
}


func main() {
    var ctx *usb.Context
    var cfgs []usb.ConfigInfo
    var keyboard usb.Device
    var ep_in usb.Endpoint
    var ep_out usb.Endpoint
    var n uint16

    flag.Parse()

    ctx = usb.NewContext()
    defer ctx.Close()
    ctx.Debug(3)

    //ctx.Debug(*debug)

    keyboard = *find_keyboard(*ctx)[0]

    cfgs = keyboard.Descriptor.Configs

    // ConfigInfo
    for _, config := range cfgs {
        // InterfaceInfo
        for _, interface_info := range config.Interfaces {
            // InterfaceSetup
            for _, interface_setup := range interface_info.Setups {
                // EndpointInfo
                for _, endpoint := range interface_setup.Endpoints {
                    fmt.Printf("Found Enpoint")
                    var err error
                    if endpoint.Direction() == usb.ENDPOINT_DIR_IN {
                        ep_in, err = keyboard.OpenEndpoint(config.Config, interface_info.Number, interface_setup.Number, endpoint.Address)
                    } else if endpoint.Direction() == usb.ENDPOINT_DIR_OUT {
                        ep_out, err = keyboard.OpenEndpoint(config.Config, interface_info.Number, interface_setup.Number, endpoint.Address)
                    }
                    if err != nil {
                        log.Fatalf("OMG!", err)
                    }
                }
            }
        }
    }

    if ep_out != nil {
        fmt.Printf("oh my  we got an output device")
    }

    n = ep_in.Info().MaxPacketSize
    b := make([]byte, n)
    for {
        time.Sleep(100 * time.Millisecond)
        get_keys(&ep_in, &b)
        fmt.Printf("%+v", b)
    }
}
