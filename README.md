# Wizhard

> A standalone golang HomeKit bridge for your Wiz bulbs

## What?

Philips Hue bulbs are pretty cool and do support HomeKit - you need a mortgage to buy them though.

On the other hand, you can get Philips Wiz bulbs (at Home Depot last I checked) for a fraction of the cost.
They are standalone (no Hue bridge), and use wifi instead of Zigbee.

The catch?
 * no HomeKit integration, so, you are stuck with Google Home or the Wiz native app
 * and... the native app is a steaming pile
 * and... google home integration is pretty terrible
 * and... if you are here, you are probably using HomeKit for everything else and wish there was a solution

So, this project lets you expose a proper HomeKit bridge (thanks to the awesome @brutella library) that control a herd of Wiz bulbs,
and let you enjoy the bliss of using HomeKit with them.

## TL;DR

```
make build
./dist/wizhard --help
./dist/wizhard register --name "Fancy fancy" --pin 87654312 --ips 1.2.3.4 --ips 5.6.7.8
```

<!--
Run the docker image, feeding it the ip addresses of your Wiz bulbs.

```
docker run -d \
    --env HomeKit_NAME="My Fancy" \
    --env HomeKit_PIN="87654312" \
    --env WIZ_IPS="1.2.3.4 5.6.7.8" \
    --name wizhard \
    --read-only \
    --cap-drop ALL \
    --net host \
    --volume /data \
    --rm \
    dubodubonduponey/HomeKit-wiz:v1
```
-->

### It works!

Cool.

* open your iPhone
* hit the (+) button (top right hand)
* then "Add Accessory"
* now "I don't have a code or cannot scan"
* hit "My_Fancy"
* "Add Anyway"
* type your pin from HomeKit_PIN above
* "Next", "Done"

### WHY... DOES... NOT... WORK

Your Wiz bulbs have to be already usable / configured on the same network.

As far as I know, there is no (mdns) discovery mechanism to figure out the ips, so you have to do that yourself.

Typically, once the bulbs are on your network, they should get a DHCP lease from your router,
so, inspect your router client table to figure it out, or nmap your way out like the grown-ups do (hint: Wiz live on port 38899).

<!--
## Roll your own, for the strong and spirited!

You need golang1.13+ to build (probably older versions work as well but can't be bothered to check).

```
make build
./wizhard --help
./wizhard register --name "Fancy fancy" --pin 87654312 --ips 1.2.3.4 --ips 5.6.7.8
```
-->
## Caveats

 * No discovery mechanism, you have to configure the bulbs ip manually after setting them up.
 * Not my fault, but yeah, the Wiz protocol is based on UDP, has no authentication, and no security whatsoever.
Not that any of these funny iot devices are secure in any way of course, but then... Wiz bulbs are just... wide open...
 * This has been hacked together quite fast, so, except bumps... see something? say something on the bugtracker - or better, submit a patch :)
 * Something funky goes on when the bulb has previously been set in one of the weird pulsating modes - in case setPilot methods fail bizarelly, consider
 a quick `echo '{"method":"getPilot","id":527,"env":"pro","params":{"mac":"0x:de:ad:be:ee:ef","rssi":-75,"cnx":"0501","src":"udp","state":true,"sceneId":0,"r":0,"g":255,"b":4,"c":0,"w":0,"dimming":100}}' | nc -u $BULBIP 38899`
 to reset it
 
