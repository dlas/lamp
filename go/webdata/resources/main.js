
function makeurl(base, map) {
	str="?"
	ka = Object.keys(map);
	for (var i = 0; i < ka.length; i++) {
		str = str + ka[i] + "=" + map[ka[i]]
		if (i != (ka.length -1)) {
			str = str +"&"
		}
	}

	return base + str
}

function setlights() {
	console.log("HELLO WORLD")
		
	rgb = {	red: $("#red").val(),
		green:$("#green").val(),
		blue:$("#blue").val()}

	
	var url = makeurl("/api/setlights", rgb);

	$.get(url)
}
