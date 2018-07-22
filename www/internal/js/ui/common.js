// debugging
//window.onerror = function(msg, url, linenumber) {
//    alert('Error message: '+msg+'\nURL: '+url+'\nLine Number: '+linenumber);
//    return true;
//};

function get_error_block_id(txt_field) {
    return txt_field + "_alert";
}

function show_field_error(txt_field, error) {
    var txt_field_selector = "#"+txt_field;
    var error_block_id = get_error_block_id(txt_field);
    var error_block_selector = "#"+error_block_id;
    $( error_block_selector ).html( "<b>"+error+"</b>" );
    $( error_block_selector ).show();

    $( txt_field_selector ).bind("keyup change", function(e) {
        $( error_block_selector ).hide();
    });
}

function hide_field_error(txt_field) {
    var error_block_id = error_block_id(txt_field);
    $( "#"+error_block_id ).hide();
}

function hide_fields_errors(form) {
    error_blocks_selector = "#"+form+" .alert";
    $( error_blocks_selector ).hide();
}

function ui_display_error(xhr, textStatus, exception) {

    if (xhr.hasOwnProperty('responseJSON')) {
        var error = xhr.responseJSON;
        if ('parameters_messages' in error) {
            for (var i = 0; i < error.parameters_messages.length; i++) {
                var pm = error.parameters_messages[i];
                var message_text = pm.messages.join('\n');
                show_field_error(pm.parameter, message_text);
            }
        } else {
            if (!('message' in error && error.message)) {
                error['message'] = 'Server Error'
            }
            $("#txt_error").text(error.message);
            $('#block_error').modal();
        }
    } else {
        window.location.href = "error.html";
    }
}