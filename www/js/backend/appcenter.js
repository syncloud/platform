$( document ).ready(function() {
    $.get( '/rest/available_apps')
        .done( function(data) {
                display_apps(data);
        }).fail( onError );
});