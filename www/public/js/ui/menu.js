$(document).ready(function() { 
	//fonts
	WebFontConfig = {
	google: { families: [ 'Roboto:400,300:latin' ] }
	};
	(function() {
	var wf = document.createElement('script');
	wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
	'://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';
	wf.type = 'text/javascript';
	wf.async = 'true';
	var s = document.getElementsByTagName('script')[0];
	s.parentNode.insertBefore(wf, s);
	})(); 
	
	// navi
	$(".menubutton").click(function(e) {
	$(".navi").toggleClass("naviopen");
	$(".menubutton").toggleClass("menuopen");
	e.preventDefault();
	});
	$(".navi a, #block1, #block2, #block3, #block4, #block5, #block6, footer").click(function(){
	$(".navi").removeClass("naviopen");
	$(".menubutton").removeClass("menuopen");
	});
});
/* When the user clicks on the button, 
	toggle between hiding and showing the dropdown content */
	function dropdown() {
	document.getElementById("myDropdown").classList.toggle("show");
	}

	// Close the dropdown menu if the user clicks outside of it
	window.onclick = function(event) {
	if (!event.target.matches('.dropdown')) {

	var dropdowns = document.getElementsByClassName("dropdown-content");
	var i;
	for (i = 0; i < dropdowns.length; i++) {
	var openDropdown = dropdowns[i];
	if (openDropdown.classList.contains('show')) {
	openDropdown.classList.remove('show');
	}
	}
	}
	}