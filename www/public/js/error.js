import * as _ from 'underscore'
import $ from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'

import '../css/site.css'
import './ui/font.js'
import './ui/common.js'
import './ui/menu.js'
import './common.js'
import './backend.static/common.js'
import './backend.static/menu.js'
import './backend.static/error.js'

function ui_send_log() {
                var btn = $("#btn_send_logs");
                btn.button('loading');
                backend.send_log(
                                function() {
                                                btn.button('reset');
                                },
                                function() {
                                                window.location.href = "index.html";
                                },
        ui_display_error
                );
}

$(document).ready(function () {
    $("#btn_send_logs").on('click', function () {
                ui_send_log();
    });
});
