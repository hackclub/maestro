# Maestro

**DEPRECATED**

_**This repository has been deprecated and is no longer maintained**_

Authors: [@TheThirdOne](https://github.com/TheThirdOne),
[@zachlatta](https://github.com/zachlatta)

-------------------------------------------------------------------------------

![](stark.gif)

-------------------------------------------------------------------------------

## Objective

The objective of Maestro is to allow a beginner to use services like 
[Twilio](https://www.twilio.com/products) through client-side javscript so they
don't have to worry about the exposed API . This will help beginners build
more/better applications which will help encourage them to continue learning.

## Usage

The JS file containing the library is found at `/static/baton.js`

Once that file is executed, the global variable `maestro` which contains all of
the module libraries will be ready. Any messages from the server that are not
caught with callbacks will automatically be logged in the console.

Open up the javascript console and [try it out!](http://maestro.ngrok.com/static/)

### Examples

```
maestro.Twilio.recieveSms("*", function(e){ //recieve any incoming sms
  maestro.Twilio.makeCall(e.From, e.To, maestro.Twilio.twiml().play(e.Body)); //call back  playing any sound file in the message
  console.log("Calling", e.From, "from", e.To, "playing", e.Body)
});
```

More examples can be found in the [docs](docs.md).

### Modules

- Twilio:
  - SMS and calling
  - Wraps both the Twilio [Restful API](https://www.twilio.com/docs/api/rest) 
  and the [TWIML API](https://www.twilio.com/docs/api/twiml)
- [Giphy](https://github.com/Giphy/GiphyAPI):
  - Image search
- Echo:
  - A simple module to help serve as template for module developers
- [Neutrino api](https://www.neutrinoapi.com/api/):
  - Currently unfinished
  - The most useful commands would be Geocode Address and Geocode Reverse
  
## License

See [COPYING](COPYING).
