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
import UiCommon from './ui/common.js'

function ui_display_device_url(data) {
		$("#txt_device_domain").attr('href', data.device_url);
		$("#txt_device_domain").text(data.device_url);
}

function ui_check_device_url() {
	backend.device_url(
   ui_display_device_url, 
   UiCommon.ui_display_error
);
}

$(document).ready(function () {

    $("#btn_reactivate").on('click', function () {
    	    backend.reactivate(
    		       function (data) {
                window.location.href = data.activate_url;
            },
            function (a, b, c) {
                UiCommon.ui_display_error(a, b, c);
            }
        );
    });

    ui_check_device_url();
});


