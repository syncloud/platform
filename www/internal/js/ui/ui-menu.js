$( document ).ready(function() {
    $("#btn_logout_large, #btn_logout_small").click(function(event) {
        backend_menu.logout({
            done: function(data) {
                window.location.href = "login.html";
            },
            fail: on_error
        });
    });

    $("#btn_restart_large, #btn_restart_small").click(function(event) {
        backend_menu.restart({
            done: function(data) {},
            fail: on_error
        });
    });

    $("#btn_shutdown_large, #btn_shutdown_small").click(function(event) {
        backend_menu.shutdown({
            done: function(data) {},
            fail: on_error
        });
    });
});
