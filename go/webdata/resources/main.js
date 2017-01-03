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
