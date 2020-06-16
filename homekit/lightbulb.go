package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/dubo-dubon-duponey/wizhard/controller"
)

type WizLightbulb struct {
	*accessory.Accessory

	Lightbulb *service.ColoredLightbulb

	Controller *controller.WizController
}

func NewWizLightbulb(address string, info accessory.Info) *WizLightbulb {
	acc := WizLightbulb{}
	acc.Accessory = accessory.New(info, accessory.TypeLightbulb)

	acc.Lightbulb = service.NewColoredLightbulb()

	acc.Controller = controller.NewWizController(address)

	acc.Lightbulb.On.OnValueRemoteUpdate(acc.Controller.SetOn)
	acc.Lightbulb.On.OnValueRemoteGet(acc.Controller.GetOn)

	acc.Lightbulb.Brightness.OnValueRemoteUpdate(acc.Controller.SetBrightness)
	acc.Lightbulb.Brightness.OnValueRemoteGet(acc.Controller.GetBrightness)

	acc.Lightbulb.Hue.OnValueRemoteUpdate(acc.Controller.SetHue)
	acc.Lightbulb.Hue.OnValueRemoteGet(acc.Controller.GetHue)

	acc.Lightbulb.Saturation.OnValueRemoteUpdate(acc.Controller.SetSaturation)
	acc.Lightbulb.Saturation.OnValueRemoteGet(acc.Controller.GetSaturation)

	acc.AddService(acc.Lightbulb.Service)

	return &acc
}
