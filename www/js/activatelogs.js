import * as URI from "uri-js";
import * as _ from 'underscore'; 
import $ from 'jquery';              
import jQuery from 'jquery';       
import 'bootstrap';             
import 'bootstrap/dist/css/bootstrap.css';   
import 'bootstrap-switch';        
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';         
import 'font-awesome/css/font-awesome.css'  
import '../css/site.css'
import '../css/material-icons.css'
import * as UiCommon from './ui/common.js'
import * as Common from './common.js'

function(values, on_always, on_done, on_error) {
    $.post("/rest/send_log", values).done(on_done).fail(on_error).always(on_always);
}

$(document).ready(function () {

    $("#form_logs").submit(function (event) {
        event.preventDefault();

        var values = $("#form_logs").serializeArray();

        var btn = $('#btn_send');
        btn.button('loading');
        $("#form_logs input").prop("disabled", true);

        send_logs(values,
            function() {
                btn.button('reset');
                $("#form_logs input").prop("disabled", false);
            },
            function(data) {
                window.location.href = "activate.html";
            },
            ui_display_error
        });
    });

});

