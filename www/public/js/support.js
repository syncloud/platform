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
import * as Common from './common.js'

function ui_display_toggles() {
	$("[type='checkbox']").each(function() {
		$(this).bootstrapSwitch();
	});
}


function ui_send_logs() {
    var tgl_support_logs = $("#tgl_support_logs");
    var btn = $("#btn_send_logs");
    btn.button('loading');

    Common.send_logs(
        tgl_support_logs.bootstrapSwitch('state'),
        function() { btn.button('reset'); },
        UiCommon.ui_display_error);
}

$(document).ready(function () {

    ui_display_toggles();
    $("#btn_send_logs").on('click', function () {
        ui_send_logs();
    });

});

