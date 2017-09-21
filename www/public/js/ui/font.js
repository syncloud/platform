$(document).ready(function() {
	//fonts
	WebFontConfig = {
	    google: { families: [ 'Roboto:400,300:latin' ] }
	};
	(function() {
        var wf = document.createElement('script');
        wf.src = 'js/ui/webfont.js';
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