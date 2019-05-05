import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'
import './ui/menu.js'
import * as UiCommon from './ui/common.js'
import * as Common from './common.js'

$(document).ready(function () {
    if (typeof mock !== 'undefined') { console.log("backend mock") };

    $("#form-login").submit(function (event) {
        event.preventDefault();

        var values = $("#form-login").serializeArray();

        var btn = $("#btn_login");
        btn.button('loading');
        $("#form-login input").prop("disabled", true);
        UiCommon.hide_fields_errors("form-login");
        $.post("/rest/login", values)
            .done((data) => {
                 window.location.href = "index.html";
            		})
            .fail(UiCommon.ui_display_error)
            .always( () => {
                 btn.button('reset');
                 $("#form-login input").prop("disabled", false);
        	    });
        
    });
});
