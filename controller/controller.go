// Package controller exposes a generic client to communicate with a single Wiz Bulb
package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dubo-dubon-duponey/wizhard/utils"
	"github.com/lucasb-eyer/go-colorful"
	"math"
)

// State represents the current or desired state of the wiz bulb
// It is being returned by a call to getPilot and should be passed as parameters to setPilot
type State struct {
	// XXX unclear what this is
	Cnx string `json:"cnx,omitempty"`

	// Bulb mac address?
	Mac string `json:"mac,omitempty"`
	// What is this? Received Signal Strength Indicator?
	Rssi int `json:"rssi,omitempty"`
	// Doesn't seem to do anything - echoed by the bulb, defaults to udp
	Src string `json:"src,omitempty"`
	// Speed - sets the color changing speed in percent - XXX not implemented for now
	Speed uint `json:"speed,omitempty"`
	// Temp - sets color temperature in kelvins - XXX not implemented for now
	Temp uint `json:"temp,omitempty"`
	// schdPsetId - rhythm id of the room - XXX not implemented for now
	SchdPsetId uint `json:"schdPsetId,omitempty"`
	/*
	   sceneId - calls one of the predefined scenes - XXX not implemented for now

	       scene 0 - Red
	       scene 0 - Green
	       scene 0 - Blue
	       scene 0 - Yellow
	       scene 1 - Ocean
	       scene 2 - Romance
	       scene 3 - Sunset
	       scene 4 - Party
	       scene 5 - Fireplace
	       scene 6 - Cozy
	       scene 7 - Forest
	       scene 8 - Pastel Colors
	       scene 9 - Wake-up
	       scene 10 - Bedtime
	       scene 11 - Warm White
	       scene 12 - Day light
	       scene 13 - Cool white
	       scene 14 - Night light
	       scene 15 - Focus
	       scene 16 - Relax
	       scene 17 - True colors
	       scene 18 - TV time
	       scene 19 - Plant growth
	       scene 20 - Spring
	       scene 21 - Summer
	       scene 22 - Fall
	       scene 23 - Deep dive
	       scene 24 - Jungle
	       scene 25 - Mojito
	       scene 26 - Club
	       scene 27 - Christmas
	       scene 28 - Halloween
	       scene 29 - Candlelight
	       scene 30 - Golden white
	       scene 31 - Pulse
	       scene 32 - Steampunk
	*/
	// XXX sceneID are problematic - if we repeat them back to the bulb while setting colors, it won't work
	// Right now, removing the scenes entirely
	// SceneId uint `json:"sceneId",omitempty`
	/*
	   State - on or off
	*/
	On bool `json:"state"`
	/*
	   r - red color range 0-255
	   g - green color range 0-255
	   b - blue color range 0-255
	   c - cold white range 0-255
	   w - warm white range 0-255
	*/
	R uint `json:"r"`
	G uint `json:"g"`
	B uint `json:"b"`
	// XXX not implemented for now
	C uint `json:"c"`
	// XXX not implemented for now
	W uint `json:"w"`
	/*
	   Bulb dimmer in percent
	*/
	Dimming uint `json:"dimming"`
}

// Method to get firmware and other system infos
const METHOD_GET_SYSTEM_CONFIG = "getSystemConfig"

// Get the bulb current state
const METHOD_GET_PILOT = "getPilot"

// Set the bulb state
const METHOD_SET_PILOT = "setPilot"

// Heartbeats sent by the bulb - XXX not implemented
const METHOD_SYNC_PILOT = "syncPilot"

// ? - XXX not implemented
const METHOD_PULSE = "Pulse"

// Register with the bulb to receive heartbeats - XXX not implemented
const METHOD_REGISTRATION = "Registration"

// QueryMessage represents a UDP message to be sent to the bulb
type QueryMessage struct {
	// Method for the messages (see method constants)
	Method string `json:"method"`
	// Not sure what is this, defaults to "pro" in all cases
	Env string `json:"env,omitempty"`
	// Possibly a unique id for the bulb?
	Id uint `json:"id,omitempty"`
	// State to pass to the bulb
	Params State `json:"params,omitempty"`
}

// Firmware represents system info returned by the bulb
type Firmware struct {
	// Mac address, presumably
	Mac string `json:"mac"`
	// Not clear what this is, seem to be static
	HomeId uint `json:"homeId"`
	// Not clear what this is
	RoomId uint `json:"roomId"`
	// Not clear what this is
	HomeLock bool `json:"homeLock"`
	// Not clear what this is
	PairingLock bool `json:"pairingLock"`
	// Not clear what this is
	TypeId uint `json:"typeId"`
	// Firmware info
	ModuleName string `json:"moduleName"`
	// Firmware info
	FwVersion string `json:"fwVersion"`
	// Not clear what this is
	GroupId uint `json:"groupId"`
	// Not clear what this is
	DrvConf []int `json:"drvConf"`
}

// ResponseSystem represents the response obtained from the bulb when querying METHOD_GET_SYSTEM_CONFIG
type ResponseSystem struct {
	Method string   `json:"method"`
	Env    string   `json:"env,omitempty"`
	Result Firmware `json:"result,omitempty"`
}

// ResponseStatus represents the response obtained from the bulb when querying METHOD_GET_PILOT
type ResponseStatus struct {
	Method string `json:"method"`
	Env    string `json:"env,omitempty"`
	State  State  `json:"result,omitempty"`
}

// An error object returned in certain cases from METHOD_SET_PILOT
type Error struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// In case of success on METHOD_SET_PILOT
type Result struct {
	Success bool `json:"success"`
}

// ReponseChange represents the bulb response to METHOD_SET_PILOT
type ResponseChange struct {
	Method string `json:"method"`
	Env    string `json:"env,omitempty"`
	Result Result `json:"result,omitempty"`
	Error  Error  `json:"error,omitempty"`
}

type WizController struct {
	Address string
	State   State
	System  Firmware
}

// Read the wiz bulb current state
func (a *WizController) Read() (err error) {
	message := QueryMessage{
		Method: "getPilot",
	}

	j, _ := json.Marshal(message)

	fmt.Println("Message we are sending:", string(j))

	response, err := utils.UDPClient(a.Address, bytes.NewReader(j))
	if err != nil {
		fmt.Println("UDP client failed dramatically", err)
		return err
	}

	fmt.Println("Response we got:", response)

	data := ResponseStatus{
		State: State{
			On: false,
			R:  1,
			G:  0,
			B:  0,
		},
	}

	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		fmt.Println("Unmarshalling response failed", response, err)
		return err
	}

	// Store the state
	a.State = data.State
	return nil
}

// Set the wiz bulb to desired state
func (a *WizController) Write() (err error) {
	// XXX deactivate the programmed scenes and shit so we do not fail dramatically
	// a.State.SceneId = 0
	a.State.Speed = 0
	a.State.C = 0
	a.State.W = 0
	a.State.Src = "udp"
	a.State.Cnx = "0501"

	message := QueryMessage{
		Method: "setPilot",
		// XXX should we use this?
		//    Id:     527,
		Env:    "pro",
		Params: a.State,
	}

	j, _ := json.Marshal(message)

	fmt.Println("Message we are sending:", string(j))

	response, err := utils.UDPClient(a.Address, bytes.NewReader(j))
	if err != nil {
		fmt.Println("UDP connection failed dramatically", err)
		return err
	}

	fmt.Println("Response we got:", response)

	data := ResponseChange{}

	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		fmt.Println("Unmarshalling response failed", response, err)
		return err
	}
	// XXX right now we don't do anything if the bulb responds with an error - implement proper error handling here
	return nil
}

// Read system and firmware information
func (a *WizController) ReadFirmwareInfo() (err error) {
	message := QueryMessage{
		Method: "getSystemConfig",
	}

	j, _ := json.Marshal(message)

	fmt.Println("Message we are sending:", string(j))

	response, err := utils.UDPClient(a.Address, bytes.NewReader(j))
	if err != nil {
		fmt.Println("UDP connection failed dramatically", err)
		return err
	}

	fmt.Println("Response we got:", response)

	data := ResponseSystem{}

	err = json.Unmarshal([]byte(response), &data)
	if err != nil {
		fmt.Println("Unmarshalling response failed", response, err)
		return err
	}
	a.System = data.Result
	return nil
}

// Initialize the controller - basically get system info and current state
func (a *WizController) Init() (err error) {
	err = a.Read()
	if err != nil {
		return err
	}
	err = a.ReadFirmwareInfo()
	if err != nil {
		return err
	}
	return nil
}

// Homekit hook to get whether the bulb is on or off
func (a *WizController) GetOn() bool {
	fmt.Println("DEBUG -> calling getOn")
	// Refresh state - not technically necessary if we are alone managing this bulb, useful if there are competing systems that changed its state after we bootstrapped
	err := a.Read()
	if err != nil {
		fmt.Println("Alas, we could not query thy noble lightbulb that appears to be dead or something")
		return false
	}
	fmt.Println("DEBUG -> answering", a.State.On)
	return a.State.On
}

// Homekit hook to set the bulb to on or off
func (a *WizController) SetOn(value bool) {
	fmt.Println("DEBUG -> calling setOn to", value)
	a.State.On = value
	err := a.Write()
	if err != nil {
		fmt.Println("Alas, we could not set thy noble lightbulb that appears to be dead or something")
	}
}

// Homekit hook to read the bulb brightness
func (a *WizController) GetBrightness() int {
	fmt.Println("DEBUG -> calling getBrightness")
	// Refresh state - not technically necessary if we are alone managing this bulb, useful if there are competing systems that changed its state after we bootstrapped
	err := a.Read()
	if err != nil {
		fmt.Println("Alas, we could not query thy noble lightbulb that appears to be dead or something")
		return 0
	}
	fmt.Println("DEBUG -> answering", a.State.Dimming)
	return int(a.State.Dimming)
}

// Homekit hook to set the bulb brightness
func (a *WizController) SetBrightness(value int) {
	fmt.Println("DEBUG -> calling setBrightness to", value)
	a.State.Dimming = uint(value)
	err := a.Write()
	if err != nil {
		fmt.Println("Alas, we could not set thy noble lightbulb that appears to be dead or something")
	}
}

// Homekit hook to read the bulb hue
func (a *WizController) GetHue() float64 {
	fmt.Println("DEBUG -> calling getHue")
	// Refresh state - not technically necessary if we are alone managing this bulb, useful if there are competing systems that changed its state after we bootstrapped
	err := a.Read()
	if err != nil {
		fmt.Println("Alas, we could not query the noble lightbulb that appears to be dead")
		return 0
	}
	color := colorful.Color{R: float64(a.State.R), G: float64(a.State.G), B: float64(a.State.B)}
	h, _, _ := color.Hsv()
	fmt.Println("DEBUG -> answering", h)
	return math.Round(h)
	//  return 0
}

// Homekit hook to set the bulb hue
func (a *WizController) SetHue(value float64) {
	fmt.Println("DEBUG -> calling setHue to", value)

	color := colorful.Color{R: float64(a.State.R), G: float64(a.State.G), B: float64(a.State.B)}
	_, s, v := color.Hsv()

	fmt.Println("Starting point", a.State.R, a.State.G, a.State.B)
	fmt.Println("WizController Set Hue", value, "with s being", s, "and v", v)

	hsv := colorful.Hsv(value, s, 255)
	fmt.Println("Hue Set Red", hsv.R)
	fmt.Println("Hue Set Green", hsv.G)
	fmt.Println("Hue Set Blue", hsv.B)

	a.State.R = uint(hsv.R)
	a.State.G = uint(hsv.G)
	a.State.B = uint(hsv.B)

	err := a.Write()
	if err != nil {
		fmt.Println("Alas, we could not set thy noble lightbulb that appears to be dead or something")
	}
}

// Homekit hook to read the bulb saturation
func (a *WizController) GetSaturation() float64 {
	fmt.Println("DEBUG -> calling getSaturation")
	// Refresh state - not technically necessary if we are alone managing this bulb, useful if there are competing systems that changed its state after we bootstrapped
	err := a.Read()
	if err != nil {
		fmt.Println("Alas, we could not query the noble lightbulb that appears to be dead")
		return 0
	}
	color := colorful.Color{R: float64(a.State.R), G: float64(a.State.G), B: float64(a.State.B)}
	_, s, _ := color.Hsv()
	fmt.Println("DEBUG -> answering", s*100)
	return math.Round(s * 100)
}

// Homekit hook to set the bulb saturation
func (a *WizController) SetSaturation(value float64) {
	fmt.Println("DEBUG -> calling setSaturation to", value)

	color := colorful.Color{R: float64(a.State.R), G: float64(a.State.G), B: float64(a.State.B)}
	h, _, v := color.Hsv()

	fmt.Println("Starting point", a.State.R, a.State.G, a.State.B)
	fmt.Println("WizController Set Saturation", value, "with h being", h, "and v", v)

	hsv := colorful.Hsv(h, value/100, 255)
	fmt.Println("Hue Set Red", hsv.R)
	fmt.Println("Hue Set Green", hsv.G)
	fmt.Println("Hue Set Blue", hsv.B)

	a.State.R = uint(hsv.R)
	a.State.G = uint(hsv.G)
	a.State.B = uint(hsv.B)

	err := a.Write()
	if err != nil {
		fmt.Println("Alas, we could not set thy noble lightbulb that appears to be dead or something")
	}
}

func NewWizController(address string) *WizController {
	wc := &WizController{
		Address: address,
		State:   State{},
	}

	// Init to get the current state in
	wc.Init()
	return wc
}

// A python implem
// https://github.com/sbidy/pywizlight

/*
Some dumps below

// echo '{"method":"setPilot","id":527,"env":"pro","params":{"mac":"0x:de:ad:be:ee:ef","rssi":-75,"cnx":"0501","src":"udp","state":true,"sceneId":0,"r":0,"g":255,"b":4,"c":0,"w":0,"dimming":100}}' | nc -u 10.0.4.208 38899
// echo '{"method":"setPilot","id":527,"env":"pro","params":{"mac":"0x:de:ad:be:ee:ef","rssi":-75,"cnx":"0501","src":"udp","state":true,"sceneId":0,"r":0,"g":255,"b":4,"c":0,"w":0,"dimming":100}}' | nc -u 10.0.4.208 38899

   params := Parameters{
     Mac:     "0x:de:ad:be:ee:ef",
     Rssi:    -75,
     Cnx:     "0501",
     Src:     "udp",
     State:   true,
     SceneId: 0,
     R:       0,
     G:       255,
     B:       0,
     C:       0,
     W:       0,
     Dimming: 50,
   }
   message := Message{
     Method: "setPilot",
     Id:     527,
     Env:    "pro",
     Params: params,
   }
*/
//    "params": {"mac":"0x:de:ad:be:ee:ef","rssi":-75,"cnx":"0501","src":"udp","state":true,"sceneId":0,"r":0,"g":255,"b":4,"c":0,"w":0,"dimming":100}
