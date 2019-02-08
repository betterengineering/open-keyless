# Open Keyless Controller
The Open Keyless Controller is designed to go on the inside of the entryway. This is to keep access to the GPIO pins on
the Raspberry Pi and the connections to the [electric door strike](strike.md) inside a secured area.

An electric door strike will need to be purchased to match your door frame and handle and is not included in this
module's parts list. For more information on choosing an electric strike, see the guide on
[Strike Buying Guide](strike.md).

## Parts List
TODO: add the housing once complete.
* [Open Keyless Controller HAT](https://oshpark.com/shared_projects/oWPmcOqR)
    * Also can be found in [assets/controller/pcbs/controller-HAT.brd](../assets/controller/pcbs/controller-HAT.brd)
    * This part will need to be manufactured. For more info, see the [overview](overview.md)
* [Breadboard-friendly 2.1mm DC barrel jack](https://www.adafruit.com/product/373)
* 1M0 / 1M ohm Resistor
    * I don't currently have a recommendation for where to pick these up. I had one on hand in my pack of resistors.
    * TODO: either redesign the controller to not need this resistor or find a reasonable supplier that wouldn't add
      another shipping fee.
* [5V 2.5A Switching Power Supply with 20AWG MicroUSB Cable](https://www.adafruit.com/product/1995)
    * This is to power the Raspberry Pi. If you already have a power source, you will not need one of these.
* [N-channel power MOSFET - 30V / 60A](https://www.adafruit.com/product/355)
* [Tapered Heat-Set Inserts for Plastic](https://www.mcmaster.com/94180a321)
    * Brass, M2.5 x 0.45 mm Thread Size, 3.4 mm Installed Length
    * This pack comes with enough for both the reader and controller. If you're building both, you only need to order
      this once.
* [Brass Heat-Set Inserts for Plastic](https://www.mcmaster.com/94459a130)
    * M3 x 0.50 mm Thread Size, 4.300 mm Installed Length
    * This pack comes with enough for both the reader and controller. If you're building both, you only need to order
      this once.
* [Stainless Steel Tamper-Resistant Button Head Torx Screws](https://www.mcmaster.com/91900a847)
    * M3 x 0.50mm Thread, 8mm Long
    * This pack comes with enough for both the reader and controller. If you're building both, you only need to order
      this once.
    * Note, if you don't have a tool for security torx screws and don't need the added security, these are a fine
      alternative: [Button Head Hex Drive Screw](https://www.mcmaster.com/92095a181)
* [White Nylon Screw and Stand-off Set – M2.5 Thread](https://www.adafruit.com/product/3658)
    * This pack comes with enough for both the reader and controller. If you're building both, you only need to order
      this once.
* [6-pin JST SM Plug + Receptacle Cable Set](https://www.adafruit.com/product/1665)
    * This pack comes with enough for both the reader and controller. If you're building both, you only need to order
      this once.
* [2-pin JST SM Plug + Receptacle Cable Set](https://www.adafruit.com/product/2880)
* [GPIO Header for Raspberry Pi HAT - 2x20 Short Female Header](https://www.adafruit.com/product/2243)
* [12V DC 1000mA (1A) regulated switching power adapter - UL listed](https://www.adafruit.com/product/798)
* [Raspberry Pi Zero W](https://www.adafruit.com/product/3400)
    * You may also need a [Mini HDMI Plug to Standard HDMI Jack Adapter](https://www.adafruit.com/product/2819) if you
      want to be able to configure the Raspberry Pi manually.
    * You may also need a [USB OTG Host Cable - MicroB OTG male to A female](https://www.adafruit.com/product/1099) if you
      want to be able to configure the Raspberry Pi manually.
* [Micro SD Card](https://www.adafruit.com/product/3259)
    * I typically buy my Micro SD cards from Amazon. I have included a link to Adafruit to limit suppliers, but any
      micro SD card over 8 GB will do.

## Assembly Guide
TODO: complete the assembly guide.
