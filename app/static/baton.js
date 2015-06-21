function Maestro(){
  var self = this;
  this.ws = new WebSocket("ws://" + window.location.host + "/baton/connect");
  var id = "";
  this.ws.onmessage = function(message){
    if(!id){
      id = message.data;
      return;
    }
    var data = JSON.parse(message.data);
    if(callbacks[data.id]){
      callbacks[data.id](data.body);
      return;
    }
    if(self[data.module]){
      self[data.module].process(data);
    }else{
      console.error("Module: " + data.module + " does not exist.");
    }
  };
  var callbacks = {};
  var counter = 0;
  this.send = function(module,call,body,callback){
    var mID = id+"-"+counter;
    if(callback){
      callbacks[mID] = callback;
    }
    this.ws.send(JSON.stringify({module:module,call:call,id:mID,body:body}));
    counter++;
  };
  this.Echo = {
    echo:function(e,callback){
      self.send("Echo","echo",e,callback);
    },
    process:function(e){
      console.log(e.body);
    }
  };
  this.Giphy = {
    search:function(e,c){
      self.send("Giphy","search",{q:e},c);
    },
    getById:function(id,c){
      if(id instanceof Array){
        self.send("Giphy","getbyids",{ids: id.join(',')},c);
      }else{
        self.send("Giphy","getbyid", {id:id},c);
      }
    },
    translate:function(e,c){
      self.send("Giphy","translate",{term:e},c);
    },
    random:function(e,c){
      self.send("Giphy","random",{tags:e},c);
    },
    trending:function(c){
      self.send("Giphy","trending","",c);
    },
    process:function(e){
      console.log(e.body);
    }
  };
  this.Neutrino = {
    process:function(e){
      console.log(e.body);
    }
  };
  this.Twilio = {
    sendSms:function(to,from,body){
      self.send("Twilio","send-sms",{to:to,from:from,body:body});
    },
    recieveSms:function(from,callback){
      self.send("Twilio","recieve-sms",{from:from},callback);
    },
    makeCall:function(to,from,twiml){
      self.send("Twilio","send-call",{to:to,from:from,twiml:twiml});
    },
    recieveCall:function(from,twiml,callback){
      self.send("Twilio","recieve-call",{from:from,twiml:twiml},callback);
    },
    process:function(e){
      console.log(e.body);
    }
  };
}
var maestro = new Maestro();