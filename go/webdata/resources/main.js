"use strict";
function makeurl(base, map) {
	var str="?"
	var ka = Object.keys(map);
	for (var i = 0; i < ka.length; i++) {
		str = str + ka[i] + "=" + map[ka[i]]
		if (i != (ka.length -1)) {
			str = str +"&"
		}
	}

	return base + str
}

function setoauth() {
	var params = {oauth: $("#oauth").val()}
	var url = (makeurl("/api/setoauth", params))
	console.log(url);
	$.get(url)
}

function setlights() {
	console.log("HELLO WORLD")
		
	var rgb = {	red: $("#red").val(),
		green:$("#green").val(),
		blue:$("#blue").val()}

	
	var url = makeurl("/api/setlights", rgb);

	$.get(url)
}

function authlink() {
	$.get("/api/getoauth?", function(data) {

		console.log(data)
		$("#authlink").attr("href", data);
	})
}

function getconfig() {
	$.get("/api/config?", function(data) {
		var p = JSON.parse(data)
		assign(p)
		console.log(p)
	})

	$.get("/api/status?", function(data) {
		var p = JSON.parse(data);
		$("#AlarmAt").val(p.AlarmAt)
		console.log(p)
	})
}

var GLOBAL_CONFIG_LIST;
function assign(data) {
	var ka = Object.keys(data);
	GLOBAL_CONFIG_LIST = ka;

	for (var i = 0; i < ka.length; i++) {
		var key = ka[i];
		var value = data[ka[i]];

		$("#" + key).val(value);
	}
}

function setconfig(data) {
	var map = rassign(GLOBAL_CONFIG_LIST);
	var url = makeurl("/api/setconfig", map);
	$.get(url, function(data) {})
}


function rassign(ids) {
	var map = {}
	for (var i = 0; i < ids.length; i++) {
		var key = ids[i];
		var value = $("#" + key).val();
		map[key] = value
	}

	return map

}

function testalarm() {
	$.get("/api/testalarm?", function(data) {})
}

window.onload =  function() {

	getconfig();
}



