Modules:
  - [Echo](#Echo)
    - [echo](#echo)
  - [Giphy](#Giphy)
    - [find](#find)
    - [findFirst](#findfirst)
    - [getById](#getbyid)
    - [translate](#translate)
    - [random](#random)
    - [trending](#trending)
  - [Neutrino](#Neutrino)
  - [Twilio](#Twilio)
    - [sendSms](#sendsms)
    - [sendMms](#sendmms)
    - [recieveSms](#recievesms)
    - [call](#call)
    - [callAndPlay](#callandplay)
    - [callAndSay](#callandsay)
    - [recieveCall](#recievecall)
    - [twiml](#twiml)
      - [say](#say)
      - [play](#play)
      - [pause](#pause)

## Echo
### echo
```
maestro.Echo.echo("Text to echo",function(reply){
  console.log(reply); //will log text to echo
});
```
## Giphy
### search
```
maestro.Giphy.search("Search text",function(response){
  console.log(response); //giphy object
});
```
### findFirst
```
maestro.Giphy.findFirst("Search text",function(response){
  console.log(response.url); //will print a url to the first result
});
```
### getById

```
maestro.Giphy.getById("FiGiRei2ICzzG",function(response){
  console.log(response); //giphy object
});
```
```
maestro.Giphy.getById(["FiGiRei2ICzzG","FiGiRei2ICzzG"],function(response){
  console.log(response); //giphy object
});
```
### translate

```
maestro.Giphy.translate("FiGiRei2ICzzG",function(response){
  console.log(response); //giphy object
});
```
### random
```
maestro.Giphy.random("dinosaur",function(response){
  console.log(response); //giphy object
});
```
### trending
```
maestro.Giphy.trending(function(response){
  console.log(response); //giphy object
});
```
##Twilio
###sendSms
```
maestro.Twilio.sendSms("human-number","twilio-number","message body");
```
###sendMms
```
maestro.Twilio.sendMms("human-number","twilio-number","http://www.vetprofessionals.com/catprofessional/images/home-cat.jpg");
```
###recieveSms
```
maestro.Twilio.recieveSms("twilio-number",function(reply){
  console.log(reply.from); //prints the number that sent a message to twilio-number
});
```
###twiml
```
var twiml = maestro.Twilio.twiml(); //create a blank twiml object

twiml.say("Hello World"); //Say Hello World
twiml.pause(3); //Pause for 3 seconds
twiml.play("http://here-and-now.info/audio/rickastley_artists.mp3"); //play "never gonna give you up"

//or alternatively
twiml.say("Hello World").pause(3).play("http://here-and-now.info/audio/rickastley_artists.mp3");
```
###call
```
maestro.Twilio.call("human-number","twilio-number",twiml);
```
###callAndPlay
```
maestro.Twilio.callAndPlay("human-number","twilio-number","http://here-and-now.info/audio/rickastley_artists.mp3");
```
###callAndSay
```
maestro.Twilio.callAndSay("human-number","twilio-number","Hello World");
```
###recieveCall
```
maestro.Twilio.recieveCall("twilio-number",twiml,function(call){
  console.log(call.from); //prints the number that called twilio-number
});
```