onValueChanged(e, t) { 
  var n = { [e]: t }; 
  if ("storagePoolName" == e && 
  (n.diskFileFormat = this.getDefaultVolumeFormat(this.props.storagePools.find(e => e.name == t))), 
  "storagePoolName" === e && this.state.mode === Ha && (n.existingVolumeName = this.getDefaultVolumeName(t)), 
  "mode" === e && t === Ha) { 
    var r = this.state.storagePoolName; 
    r && (n.existingVolumeName = this.getDefaultVolumeName(r)) 
  } this.setState(n) 
}

onValueChanged(e, t) { 
  var n = { [e]: t }; 
  if (this.setState(n), "networkType" == e) { 
    var r = !1; 
    if ("network" == t) { 
      var i = this.props.networks.map(e => e.name); i.length > 0 ? this.setState({ networkSource: i[0] }) : (this.setState({ networkSource: void 0 }), r = !0) 
    } this.setState({ saveDisabled: r }) 
  } 
}

onValueChanged(e, t) { 
  var n = {}; 
  if ("memory" == e) { 
    var r = Object(g.c)(t, this.state.memoryUnit, "KiB"); 
    r <= this.state.maxMemory ? n.memory = r : r > this.state.maxMemory && "running" != this.props.vm.state && (n.memory = r, n.maxMemory = r) 
  } else if ("maxMemory" == e) { 
    var i = Object(g.c)(t, this.state.maxMemoryUnit, "KiB"); 
    i < this.state.nodeMaxMemory && (n.maxMemory = i) 
  } else "memoryUnit" != e && "maxMemoryUnit" != e || (n = { [e]: t }); 
  this.setState(n) 
}

onValueChanged(e, t) { 
  var n = { [e]: t }; 
  this.setState(n)
}

onValueChanged(e, t) { 
  if ("source" == e) { 
    var n = Object.keys(t)[0], r = t[Object.keys(t)[0]]; 
    this.setState({ source: Object.assign({}, this.state.source, { [n]: r }) }) 
  } else "type" == e ? 
  ("disk" == t ? this.setState({ source: Object.assign({}, this.state.source, { format: "dos" }) }) : 
  this.setState({ source: Object.assign({}, this.state.source, { format: void 0 }) }), this.setState({ [e]: t })) : this.setState({ [e]: t }) 
}