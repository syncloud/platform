$( document ).ready(function() {
    $("#btn_logout_large, #btn_logout_small").click(function(event) {
        var posting = $.post( '/rest/logout' )
            .done( function(data) {
                window.location.href = "login.html";
            })
            .fail( function(xhr, textStatus, errorThrown) {
            });
    });
});
