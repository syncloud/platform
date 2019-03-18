import * as _ from 'underscore'; 
import $ from 'jquery';
import jQuery from 'jquery';

import 'bootstrap';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-switch';
import 'bootstrap-switch/dist/css/bootstrap3/bootstrap-switch.css';
import 'font-awesome/css/font-awesome.css'
import '../css/site.css'

import './ui/font.js'
import UiCommon from './ui/common.js'
import './ui/menu.js'

import * as Common from './common.js'

import Templates from './storage.templates.js'

function update_disks(on_complete, on_error) {
    $.get('/rest/settings/disks').done(on_complete).fail(on_error);
}

function update_boot_disk(on_complete, on_error) {
    $.get('/rest/settings/boot_disk').done(on_complete).fail(on_error);
}

export function boot_extend(on_complete, on_error) {
    $.get('/rest/settings/boot_extend')
        .done(function (data) {
                      Common.check_for_service_error(data, function () {
                          Common.run_after_job_is_complete(
                              Common.job_status,
                              setTimeout,
                              function () {
                                  update_boot_disk(
                                      on_complete,
                                      on_error);
                              }, on_error, 'boot_extend');
                      }, on_error);
                  })
        .fail(on_error);

}

export function disk_action(disk_device, is_activate, on_always, on_error) {
    var mode = is_activate ? "disk_activate" : "disk_deactivate";
    $.get('/rest/settings/' + mode, {device: disk_device})
        .done(function(data) {
              Common.check_for_service_error(
                  data,
                  function() {},
                  on_error);
          })
        .always(on_always)
        .fail(on_error);
}


export function disk_format(disk_device, on_complete, on_error) {
    $.post('/rest/storage/disk_format', {device: disk_device})
        .done(function (data) {
                      Common.check_for_service_error(data, function () {
                          Common.run_after_job_is_complete(
                              Common.job_status,
                              setTimeout,
                              on_complete,
                              on_error,
                              'disk_format');
                      }, on_error);
                  })
        .fail(on_error);
}

function ui_display_toggles() {
	$("[type='checkbox']").each(function() {
		$(this).bootstrapSwitch();
	});
}

function ui_enable_controls(enabled) {
	$("[type='checkbox']").each(function() {
	    $(this).bootstrapSwitch('disabled', !enabled);
	});

	$("[data-type='format']").each(function() {
		$(this).prop('disabled', !enabled);
	});

}

function ui_display_disks(data) {
		$("#block_disks").html(_.template(Templates.Disks)(data));
		ui_display_toggles();
		ui_enable_controls(true);

		$("#block_disks").find("[type='checkbox']").each(function() {
			var tgl = $(this);
			tgl.off('switchChange.bootstrapSwitch').on('switchChange.bootstrapSwitch', function(e, s) {
				var state = tgl.bootstrapSwitch('state');

                $('#partition_state').val(state);
                $('#partition_device').val(tgl.data('partition-device'));

			    var partition_index = tgl.data('partition-index');
			    var disk_index = tgl.data('disk-index');
                $('#partition_id').val(disk_index + '_' + partition_index);

                $('#partition_name').html($('#partition_name_' + disk_index + '_' + partition_index).html());
                $('#partition_disk_name').html($('#disk_name_' + disk_index).html());

                $('#partition_action_confirmation').modal('show');

			});
		});

  $("#btn_partition_action").off('click').on('click', function(e, s) {

    var state = $('#partition_state').val() == 'true';
    var tgl_loading = $('#tgl_partition_' + $('#partition_id').val() + '_loading');
    $(tgl_loading).addClass('opacity-visible');

    ui_enable_controls(false);

    var device = $('#partition_device').val();
    disk_action(device, state, ui_check_disks, UiCommon.ui_display_error);
  });
  
  $('#partition_action_confirmation').on('hidden.bs.modal', function () {
    var state = $('#partition_state').val() == 'true';
    var tgl = $('#tgl_partition_' + $('#partition_id').val());
    tgl.bootstrapSwitch('state', !state, true);
  });

		$("#block_disks").find("[data-type='format']").each(function() {
			var btn = $(this);
			btn.button('reset');
			btn.off('click').on('click', function(e, s) {
			    var index = btn.data('index');
                $('#disk_index').val(index);
                $('#disk_name').html($('#disk_name_' + index).html());
                $('#disk_format_confirmation').modal('show');
            });
		});

        $("#btn_disk_format").off('click').on('click', function(e, s) {

            var index = $('#disk_index').val();
            var btn = $('#btn_format_' + index);
            btn.button('loading');
            ui_enable_controls(false);

            var device = btn.data('device');

            disk_format(device,
                ui_check_disks,
                function (a, b, c) {
                    UiCommon.ui_display_error(a, b, c);
                    btn.button('reset');
                });
        });

}

function ui_display_boot_disk(data) {
		$("#block_boot_disk").html(_.template(Templates.BootDisk)(data.data));
		var btn = $("#btn_boot_extend");
        btn.button('reset');
        btn.off('click').on('click', function () {
    		ui_boot_extend();
        });
}

function ui_boot_extend() {
    var btn = $("#btn_boot_extend");
    btn.button('loading');
    
    boot_extend(
        ui_display_boot_disk,
        function (a, b, c) {
            UiCommon.ui_display_error(a, b, c);
            btn.button('reset');
        });

}

function ui_check_disks() {
		update_disks(ui_display_disks, UiCommon.ui_display_error);
}

function ui_check_boot_disk() {
		update_boot_disk(ui_display_boot_disk, UiCommon.ui_display_error);
}

$(document).ready(function () {
    if (typeof mock !== 'undefined') { console.log("backend mock") };
    ui_check_disks();
    ui_check_boot_disk();

});
