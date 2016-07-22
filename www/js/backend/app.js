function register_btn_action_click(app_id, action) {
		$("#btn_"+action).on('click', function () {
				var btn = $(this);
				btn.button('loading');
				$.get("/rest/"+action, {app_id: app_id})
						.always(function() {
								run_after_sam_is_complete(function() {
										update_app(app_id, function() {
												btn.button('reset');
										});
								});
						});

		});
}

function update_app(app_id, on_complete) {
		$.get( '/rest/app', {app_id: app_id})
						.done( function(data) {
								display_app(data);
								$("#btn_open").on('click', function () {
										var btn = $(this);
										var app_url = btn.data('url');
        						window.location.href = app_url;
								});
								register_btn_action_click(app_id, 'install');
								register_btn_action_click(app_id, 'upgrade');
								register_btn_action_click(app_id, 'remove');
						})
            .always(function() {
            		typeof on_complete === 'function' && on_complete();
            });
}

$( document ).ready(function () {
		var app_id = new URI().query(true)['app_id'];

		update_app(app_id);
});
