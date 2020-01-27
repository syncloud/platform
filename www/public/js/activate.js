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


function activate(parameters, on_always, on_done, on_error) {
        $.post("/rest/activate", parameters)
            .done(on_done)
            .fail(on_error)
            .always(on_always);
}

function activate_custom_domain(parameters, on_always, on_done, on_error) {
    $.post("/rest/activate_custom_domain", parameters)
            .done(on_done)
            .fail(on_error)
            .always(on_always);
}

function login() {
        var url = (new URI())
                .protocol('https')
                .port(443)
                .filename("")
                .query("");

        window.location.href = url;
}


$(document).ready(function () {

    if (typeof mock !== 'undefined') { console.log("backend mock") };

        $("#domain_type_syncloud").click(function (event) {
				event.preventDefault();
				$("#domain_type").val('syncloud');
		});

		$("#domain_type_custom").click(function (event) {
				event.preventDefault();
				$("#domain_type").val('custom');
		});

    $("#btn_error_send_logs").on('click', function () {
        window.location.href = "sendlogs.html";
    });

		$("#form_activate").submit(function (event) {
				event.preventDefault();

				var values = $("#form_activate").serializeArray();

				var btn = $('#btn_activate');
				btn.button('loading');
				$("#form_activate input").prop("disabled", true);
				hide_fields_errors("form_activate");

				var on_always = function() {
					btn.button('reset');
					$("#form_activate input").prop("disabled", false);
				};

   				if ( $("#domain_type").val() == 'syncloud') {
					activate(
						values,
						on_always,
						backend.login,
						ui_display_error);
				} else {
					activate_custom_domain(
						values,
						on_always,
						backend.login,
						ui_display_error);
				}
		});
});

