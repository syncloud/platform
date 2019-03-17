import URI from "urijs";
import * as _ from 'underscore';
import $ from 'jquery';
import jQuery from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'
import { AppTemplate } from './app.templates.js'
import './ui/font.js'
import UiCommon from './ui/common.js'
import './ui/menu.js'

import * as Common from './common.js'
import './backend/menu.js'

function load_app(app_id, on_complete, on_error) {
    $.get('/rest/app', {app_id: app_id}).done(on_complete).fail(on_error);
};

function app_action(app_id, action, on_always, on_error) {
    $.get("/rest/" + action, {app_id: app_id}).always(on_always).fail(on_error);
};

export function run_app_action(app_id, action, on_complete, on_error) {
    app_action(app_id, action, function (data) {
        Common.check_for_service_error(data, function () {
            Common.run_after_sam_is_complete(
                Common.job_status,
                setTimeout,
                on_complete,
                on_error);
        }, on_error)
    }, on_error);
}

function register_btn_open_click() {
    $("#btn_open").off('click').on('click', function () {
        var btn = $(this);
        var app_url = btn.data('url');
        window.location.href = app_url;
    });
}

function register_btn_action_click(app_id, action) {

    $("#btn_" + action).off('click').on('click', function () {
         $('#app_id').val(app_id);
         $('#app_action').val(action);
         $('#confirm_caption').html($('#btn_' + action).html());
         $('#app_action_confirmation').modal('show');
    });
}

function app_action(app_id, action) {

        var btn = $("#btn_" + action);
        btn.button('loading');

        run_app_action(app_id, action, function () {
            btn.button('reset');
            ui_load_app();
        }, UiCommon.ui_display_error);

}

function ui_display_app(data) {
		$("#block_app").html(_.template(AppTemplate)(data));
		var app_id = data.info.app.id;
		register_btn_open_click();
		register_btn_action_click(app_id, 'install');
		register_btn_action_click(app_id, 'upgrade');
		register_btn_action_click(app_id, 'remove');
		
	    $("#btn_confirm").off('click').on('click', function () {
		    var app_id = $('#app_id').val();
            var action =  $('#app_action').val();
            app_action(app_id, action);
         });
}

function ui_load_app() {
		var app_id = new URI().query(true)['app_id'];

		load_app(app_id, ui_display_app, UiCommon.ui_display_error);
}

$( document ).ready(function () {
		ui_load_app();
});
