function Maestro(){
  var self = this;
  this.ws = new WebSocket("ws://localhost:1759/baton/connect")
  this.ws.onmessage = function(message){
    var data = JSON.parse(message.data);
    if(self[data.module]){
      self[data.module].process(data);
    }else{
      console.error("Module: " + data.module + " does not exist.");
    }
  };
  this.send = function(module,call,body){
    this.ws.send(JSON.stringify({module:module,call:call,body:body}));
  };
  this.Echo = {
    echo:function(e){
      self.send("Echo","echo",e);
    },
    process:function(e){
      console.log(e.body);
    }
  };
  this.Giphy = {
    search:function(e){
      self.send("Giphy","search",{q:e});
    },
    getById:function(id){
      if(id instanceof Array){
        self.send("Giphy","getbyids",{ids: id.join(',')});
      }else{
        self.send("Giphy","getbyid", {id:id});
      }
    },
    translate:function(e){
      self.send("Giphy","translate",{term:e});
    },
    random:function(e){
      self.send("Giphy","random",{tags:e});
    },
    trending:function(){
      self.send("Giphy","trending","");
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
}
var maestro = new Maestro();