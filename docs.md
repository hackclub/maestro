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
### find
```
maestro.Giphy.find("Search text",function(response){
  console.log(response);
})
```
### findFirst
### getById
### translate
### random
### trending
##Twilio
###sendSms
###sendMms
###recieveSms
###call
###callAndPlay
###callAndSay
###recieveCall
###twiml
####say
####play
####pause
