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
});