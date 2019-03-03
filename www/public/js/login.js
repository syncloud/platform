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
import './ui/menu.js'
import './backend/common.js'
import './backend/menu.js'
import './backend/login.js'

$(document).ready(function () {
    $("#form-login").submit(function (event) {
        event.preventDefault();

        var values = $("#form-login").serializeArray();

        var btn = $("#btn_login");
        btn.button('loading');
        $("#form-login input").prop("disabled", true);
        UiCommon.hide_fields_errors("form-login");

        backend.login(
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
