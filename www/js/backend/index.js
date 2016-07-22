$( document ).ready(function() {
    $.get( '/rest/installed_apps')
        .done( function(data) {
                display_apps(data);
        }).fail( onError );
});
