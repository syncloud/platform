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
        backend_menu.logout({
            done: function(data) {
                window.location.href = "login.html";
            },
            fail: ui_display_error
        });
    });

    $("#btn_restart_large, #btn_restart_small").click(function(event) {
        backend_menu.restart({
            done: function(data) {},
            fail: ui_display_error
        });
    });

    $("#btn_shutdown_large, #btn_shutdown_small").click(function(event) {
        backend_menu.shutdown({
            done: function(data) {},
            fail: ui_display_error
        });
    });
});
