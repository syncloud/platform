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
import './ui/menu.js'

function device_url(on_complete, on_error) {
        $.get('/rest/settings/device_url').done(on_complete).fail(on_error);
    };

function reactivate(on_complete, on_error) {
        $.post('/rest/settings/deactivate').done(on_complete).fail(on_error);
    };

function ui_display_device_url(data) {
		$("#txt_device_domain").attr('href', data.device_url);
		$("#txt_device_domain").text(data.device_url);
}

function ui_check_device_url() {
	device_url(
        ui_display_device_url,
        UiCommon.ui_display_error
    );
}

$(document).ready(function () {
    if (typeof mock !== 'undefined') { console.log("backend mock") };

    $("#btn_reactivate").on('click', function () {
    	reactivate(
    		function (data) {
                window.location.href = "/";
            },
            function (a, b, c) {
                UiCommon.ui_display_error(a, b, c);
            }
        );
    });

    ui_check_device_url();
});


