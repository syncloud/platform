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
import UiCommon from './ui/common.js'
import Common from './common.js'

function login(values, on_complete, on_error, on_always) {
        $.post("/rest/login", values).done(on_complete).fail(on_error).always(on_always);
    };

$(document).ready(function () {
    $("#form-login").submit(function (event) {
        event.preventDefault();

        var values = $("#form-login").serializeArray();

        var btn = $("#btn_login");
        btn.button('loading');
        $("#form-login input").prop("disabled", true);
        UiCommon.hide_fields_errors("form-login");

        login(
                values,
        		function(data) {
                    window.location.href = "index.html";
        		},
                UiCommon.ui_display_error,
                function() {
                    btn.button('reset');
                    $("#form-login input").prop("disabled", false);
        		}
        );
    });
});
