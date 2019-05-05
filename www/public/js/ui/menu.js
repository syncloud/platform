import * as UiCommon from './common.js'
import styles from 'roboto-fontface/css/roboto/roboto-fontface.css'

function logout(on_complete, on_error) {
    $.post('/rest/logout').done(on_complete).fail(on_error);
}

function restart(on_complete, on_error) {
    $.get('/rest/restart').done(on_complete).fail(on_error);
}

function shutdown(on_complete, on_error) {
    $.get('/rest/shutdown').done(on_complete).fail(on_error);
}


// When the user clicks on the button, toggle between hiding and showing the dropdown content
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

$( document ).ready(function() {

	$(".menubutton").click(function(e) {
        $(".navi").toggleClass("naviopen");
        $(".menubutton").toggleClass("menuopen");
        e.preventDefault();
	});
	$(".navi a, #block1, #block2, #block3, #block4, #block5, #block6, footer").click(function(){
	    $(".navi").removeClass("naviopen");
	    $(".menubutton").removeClass("menuopen");
	});


    $("#btn_logout_large, #btn_logout_small").click(function(event) {

        logout(
            function(data) {
                window.location.href = "login.html";
            }, 
            UiCommon.ui_display_error);
    });

    $("#btn_restart_large, #btn_restart_small").click(function(event) {
        restart(
            function(data) { },
            UiCommon.ui_display_error
        );
    });

    $("#btn_shutdown_large, #btn_shutdown_small").click(function(event) {
        shutdown(
            function(data) { },
            UiCommon.ui_display_error
        );
    });
});
