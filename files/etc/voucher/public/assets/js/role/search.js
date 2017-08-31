$(document).ready(function () {
	getRole();
});

function getRole() {
	$.ajax({
		url: '/v1/ui/role/all',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			for (i = 0; i < arrData.length; i++) {
				var html = "<div class='card'><div class='card-body pt0 pb0'><div class='row'>"
					+ "<div class='col-sm-9'>"
					+ "<div class='checkbox c-checkbox'>"
					+ "<label>"
					+ arrData[i].detail
					+ "</label>"
					+ "</div>"
					+ "</div>"
					+ "<div class='col-sm-3'>"
					+ "<button type='button' onclick='edit(\"" + arrData[i].id + "\")' class='btn btn-raised btn-sm btn-info down-5px'><em class='ion-edit'></em></button>"
					+ "</div>"
					+ "</div></div></div>";
				var li = $("<div class='col-md-3'></div>").html(html);
				li.appendTo('#listTag');
			}
		}
	});
}

function add(param) {
	window, location = "/role/create";
}

function edit(param) {
	window.location = "/role/edit?id=" + param;
}
