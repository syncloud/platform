Disks = `
    <%
		if (disks.length == 0) {
    %>

		<span class="span">No external disks found</span>

    <% } %>

    <%
    for (i=0; i < disks.length; i++) {
				var disk = disks[i];
    %>

		<div class="setline" style="margin-top: 20px;">
			<span class="span" style="font-weight: bold;" id="disk_name_<%=i%>"><%= disk.name %> - <%= disk.size %></span>
			<% if (!disk.active) { %>
			<div class="spandiv">
				<button class="buttonred bwidth smbutton btn-lg"
						id="btn_format_<%=i%>"
						data-type="format"
                        data-index="<%=i%>"
                        data-device="<%=disk.device %>"
						data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> ">Format</button>
			</div>
			<% } %>
    	</div>

    <%
				partitions = disk.partitions;
				for (j=0; j < partitions.length; j++) {
				var partition = partitions[j];

					if (partition.mountable || partition.active) {
    %>

		<div class="setline">
				<span class="span" id="partition_name_<%=i%>_<%=j%>">Partition - <%=partition.size %></span>
				<div class="spandiv">
						<input type="checkbox" id="tgl_partition_<%=i%>_<%=j%>"
							   data-disk-index="<%=i%>"
                               data-partition-index="<%=j%>"
                               data-partition-device="<%=partition.device %>"
							   data-on-text="Active"
							   data-off-text="Not active"
							   data-label-width="8" <% if (partition.active) { %>checked<% } %> />
						<i class="fa fa-circle-o-notch fa-spin switchloading opacity-invisible"
						   id="tgl_partition_<%=i%>_<%=j%>_loading"></i>
				</div>
		</div>

    <%
					}

            	}
    %>

    <% } %>
`

BootDisk = `
        <span class="span">Partition - <%=size %></span>
        <% if (extendable) { %>
        <div class="spandiv">
            <button class="buttongreen bwidth smbutton btn-lg"
					id="btn_boot_extend"
					data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i> Extending...">Extend</button>
        </div>
        <% } %>
`

module.exports = {
	Disks,
	BootDisk
};