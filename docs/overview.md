# Overview
This project contains two modules: a [reader](reader.md) and a [controller](controller.md). The reader module is meant
to sit on the outside of the door to scan RFID badges and the controller sits on the inside of the door to control the
electric door strike and determine if the user has access.

At the time of writing, this system is a one to one mapping in that one controller can only service one reader. This
project will likely be extended in the next iteration to support more reader/controller combinations and may require
a centralized controller to allow for many badge readers.

## Cost
At the time of writing, I calculated the cost for building the reader and controller using all links provided to be
about $200 USD ($208.55 to be precise) including shipping. Your end cost may vary. This project was not optimized for
cost in the current iteration, so this could easily be optimized to be lower.

This cost does not include the electric strike. The electric strike will likely be the most expensive part of this
build. I've seen electric strikes vary from $50 to $400 USD depending on the grade and use case of the strike, so keep
that in mind.

If this project acquires enough interest, I would be willing to sell assemble it yourself kits with a slight markup that
would lower the cost of overall because I could buy in bulk.

## Printed Circuit Boards
This project contains several printed circuit boards (PCBs) that need to be manufactured. These boards were designed
using [Eagle](https://www.autodesk.com/products/eagle/overview) and can be found in the `/assets` directory in this
repository. If you don't have a PCB manufacturer of choice, I recommend using [OSH Park](https://oshpark.com/). They are
well made and inexpensive. They do have a minimum order of 3 units, but I have found the final cost to be so low
($2-$4 USD per board and free shipping) that 3 units is still affordable even if you only need one.

In each module guide, I have included a shared project link that you can order directly from OSH Park without needing to
upload the `.brd` file.

## 3D Printed Housing
This project contains several 3D models that need to be manufactured. These models were designed via
[Onshape](https://www.onshape.com/) and can be found publicly
[here](https://cad.onshape.com/documents/67f068fc5317736e5689dda5/w/9dc79127656737c6368d36b9/e/1df9f27c5d40799b76486afd).
Additionally, I have included the models in SolidWorks and Parasolid formats in the `/assets` directory of this repo 
for convenience. If you don't have a 3D printing service of choice, I recommend
[Voodoo Manufacturing](https://voodoomfg.com/). They offer reasonable prices for 3D printing services.

Unfortunately, Voodoo does not currently have the functionality to allow me to host the project for ordering, so you
will need to upload the `.stl` files to Voodoo in order to order the parts needed.

## Components
I have a preference for [Adafruit](https://www.adafruit.com/) for electronics components and
[McMaster-Carr](https://www.mcmaster.com/) for physical components. All components linked in this guide should give you
enough details to swap out components from other suppliers or from components you have on hand and there is no
requirement to order from either supplier.

## See Also
* [Reader Overview](reader.md) - An overview and guide for the Open Keyless Reader module
* [Contorller Overview](controller.md) - An overview and guide for the Open Keyless Controller module
