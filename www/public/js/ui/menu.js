import UiCommon from './common.js'

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
