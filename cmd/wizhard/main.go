package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/dubo-dubon-duponey/wizhard/homekit"
	"github.com/dubo-dubon-duponey/wizhard/utils"
	"github.com/urfave/cli"
	"log"
	"os"
)

var bulb *homekit.WizLightbulb

func register(c *cli.Context) error {
	ips := c.StringSlice("ips")
	pin := c.String("pin")
	storage := c.String("data-path")
	port := c.String("port")

	info := accessory.Info{
		Name:             c.String("name"),
		Manufacturer:     c.String("manufacturer"),
		SerialNumber:     c.String("serial"),
		Model:            c.String("model"),
		FirmwareRevision: c.String("version"),
	}

	if len(ips) == 0 {
		fmt.Println("Hey! You need to provide at least one ip! These bulbs are not going to get to work on themselves!")
	}

	//  ip := fmt.Sprintf("%s:38899", ips[0])
	//  bulb := homekit.NewWizLightbulb(ip, info)

	bulbs := []*accessory.Accessory{}

	for x, ip := range ips {
		fmt.Println("Addr:", ip)
		ip = fmt.Sprintf("%s:38899", ip)
		u, _ := utils.GenerateUUID()
		n := fmt.Sprintf("Wiz %d", x)
		fmt.Println("Bulb info", n, ip)
		bulb = homekit.NewWizLightbulb(ip, accessory.Info{
			Name:             n,
			Manufacturer:     info.Manufacturer,
			SerialNumber:     u,
			Model:            "Bulby",
			FirmwareRevision: info.FirmwareRevision,
		})
		bulbs = append(bulbs, bulb.Accessory)
	}

	bridge := accessory.NewBridge(info)

	t, err := hc.NewIPTransport(hc.Config{
		Pin:         pin,
		StoragePath: storage,
		Port:        port,
	}, bridge.Accessory, bulbs...) // bulb.Accessory)

	/*  t, err := hc.NewIPTransport(hc.Config{
	      Pin:         pin,
	      StoragePath: storage,
	      Port:        port,
	    }, bridge.Accessory, bulbs...) // bulb.Accessory)
	*/
	if err != nil {
		return err
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()

	return nil
}

/*
info := accessory.Info{
Name: "WizLamp",
SerialNumber: "1234-23AAM1",
Manufacturer: "DuboDubonDuponey",
Model: "A",
FirmwareRevision: "0.0.1",
}

ac := homekit.NewWizLightbulb("10.0.4.208:38899", info)

// configure the ip transport
config := hc.Config{Pin: "14041976"}
t, err := hc.NewIPTransport(config, ac.Accessory)
if err != nil {
log.Panic(err)
}

hc.OnTermination(func(){
  <-t.Stop()
})

t.Start()
*/

func main() {

	app := cli.NewApp()
	app.Name = "HomeKit WizHard"
	app.Usage = "Control your Wiz bulbs over HomeKit"

	uuid, _ := utils.GenerateUUID()

	app.Commands = []cli.Command{
		{
			Name:   "register",
			Usage:  "register a HomeKit device",
			Action: register,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "pin",
					Value: "87654312",
					Usage: "Pin code for your device (8 characters)",
				},
				cli.StringFlag{
					Name:  "port",
					Value: "12345",
					Usage: "Port to expose the service on",
				},
				cli.StringFlag{
					Name:  "name",
					Value: "Dubo Dubon Duponey WizHard",
					Usage: "Name of your Wiz bridge",
				},
				cli.StringFlag{
					Name:  "data-path",
					Value: "/tmp/dubo-wizhard",
					Usage: "Where to store the data files for that device",
				},
				cli.StringFlag{
					Name:  "manufacturer",
					Value: "Dubo Dubon Duponey",
					Usage: "Manufacturer of your bridge",
				},
				cli.StringFlag{
					Name:  "serial",
					Value: uuid,
					Usage: "Serial number of your bridge",
				},
				cli.StringFlag{
					Name:  "model",
					Value: "WizHard Bridge",
					Usage: "Model of your bridge",
				},
				cli.StringFlag{
					Name:  "version",
					Value: "1",
					Usage: "Firmware version of your bridge",
				},
				cli.StringSliceFlag{
					Name:  "ips",
					Usage: "IPs addresses of your bulbs",
				},
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

}

/*
  type WizBulb struct {
    *service.Service

    On *characteristic.On
    Brightness *characteristic.Brightness
    Hue *characteristic.Hue
    Saturation *characteristic.Saturation
  }

  const TypeLightbulb = "43"

  func NewLightbulb() *WizBulb {
    svc := WizBulb{}
    svc.Service = service.New(TypeLightbulb)

    svc.On = characteristic.NewOn()
    svc.AddCharacteristic(svc.On.Characteristic)

    svc.Brightness = characteristic.NewBrightness()
    svc.AddCharacteristic(svc.Brightness.Characteristic)

    svc.Hue = characteristic.NewHue()
    svc.AddCharacteristic(svc.Hue.Characteristic)

    svc.Saturation = characteristic.NewSaturation()
    svc.AddCharacteristic(svc.Saturation.Characteristic)

    return &svc
  }
*/

/*
  svc.On.OnValueRemoteUpdate(sosy.Controller.SetOn)
  svc.On.OnValueRemoteGet(sosy.Controller.GetOn)

  sosy.HomeKit.Amplifier.Volume.OnValueRemoteUpdate(sosy.Controller.SetVolume)
  sosy.HomeKit.Amplifier.Volume.OnValueRemoteGet(sosy.Controller.GetVolume)
*/

func mainOld() {
	// create an accessory
	info := accessory.Info{
		Name:             "WizLamp",
		SerialNumber:     "1234-23AAM1",
		Manufacturer:     "DuboDubonDuponey",
		Model:            "A",
		FirmwareRevision: "0.0.1",
	}

	ac := homekit.NewWizLightbulb("10.0.4.208:38899", info)

	// configure the ip transport
	config := hc.Config{Pin: "14041976"}
	t, err := hc.NewIPTransport(config, ac.Accessory)
	if err != nil {
		log.Panic(err)
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}

/*
bridge := accessory.NewBridge(...)
outlet := accessory.NewOutlet(...)
lightbulb := accessory.NewColoredLightbulb(...)

hc.NewIPTransport(config, bridge, outlet.Accessory, lightbulb.Accessory)
*/
