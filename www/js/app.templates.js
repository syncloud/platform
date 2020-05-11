const AppTemplate = `
<div class="appblock">
			<div>
				<div>
					<img src="<%= info.app.icon %>" class="appimg" alt="">
				</div>

				<div class="appinfo">
				<h1><%= info.app.name %></h1>




<div id="app_info" class="modal fade bs-are-use-sure" tabindex="-1" role="dialog" aria-labelledby="mySmallModalLabel">
	<div class="modal-dialog" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
				<h4 class="modal-title">App information</h4>
			</div>
			<div class="modal-body">
				<div class="bodymod">
					<div class="btext">
					<h2><%= info.app.name %></h2>
					<div class="btext"><%= info.app.description %></div>
			<h2>Details</h2>

        	<% if (info.installed_version !== null) { %>
					<b>Installed version:</b> <%= info.installed_version %><br>
        	<% } %>
        	<% if (info.installed_version !== info.current_version) { %>
					<b>Available version:</b> <%= info.current_version %><br>
        	<% } %>
					<!--<b>Size:</b> 17.0 MB-->
					<!--<b>Website:</b> <a href="" target="_blank">content</a>-->
					<!--<b>Licence:</b> content-->
					<!--<b>Last update:</b> content-->
					<!--<b>Channel:</b> content-->
					
					</div>

				</div>
				<div class="modal-footer">
					<button type="button" class="btn buttonlight bwidth smbutton" data-dismiss="modal">Close</button>
				</div>
			</div>
		</div>
	</div>
</div>





				</div>
			</div>
			<div>
				<div class="buttonblock">
        	<% if (info.installed_version !== null) { %>
					<button id="btn_open" data-url="<%= info.app.url %>" class="buttonblue bwidth smbutton">Open</button>
        	<% } %>
        	<% if (info.installed_version === null) { %>
					<button id="btn_install" class="buttonblue bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Installing...">Install</button>
        	<% } %>
        	<% if (info.installed_version !== null && info.installed_version !== info.current_version) { %>
					<button id="btn_upgrade" class="buttongreen bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Upgrading...">Upgrade</button>
        	<% } %>
        	<% if (info.installed_version !== null) { %>
					<button id="btn_remove"  class="buttongrey bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Removing...">Remove</button>
        	<% } %>
        <% if (info.installed_version !== null) { %>
					<button id="btn_backup"  class="buttonblue bwidth smbutton" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Creating backup...">Backup</button>
        	<% } %>
        	<button data-toggle="modal" data-target="#app_info" type="button" class="control" style="background:transparent;">
				<i class='fa fa-info-circle fa-lg'></i>
			</button>
				</div>
			</div>
		</div>
`

module.exports = {
	AppTemplate
};
