import styles from 'roboto-fontface/css/roboto-condensed/roboto-condensed-fontface.css'

$(document).ready(function() {
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
