import * as _ from 'underscore'
import $ from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/material-icons.css'

import '../css/site.css'
import './ui/common.js'
import './ui/menu.js'
import * as Common from './common.js'

function ui_send_log() {
    var btn = $("#btn_send_logs");
    btn.button('loading');
    Common.send_log(
        function() { btn.button('reset'); },
        function() { window.location.href = "index.html"; },
        ui_display_error
    );
}

$(document).ready(function () {
    $("#btn_send_logs").on('click', function () {
                ui_send_log();
    });
});
