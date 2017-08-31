$(document).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	getTag();
});

function getTag() {
	$.ajax({
		url: '/v1/ui/tag/all',
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
					+ "<input name='agreements' value='" + arrData[i] + "' type='checkbox'><span class='ion-checkmark-round'></span>" + arrData[i]
					+ "</label>"
					+ "</div>"
					+ "</div>"
					+ "<div class='col-sm-3'>"
					+ "<button type='button' value='" + arrData[i] + "' class='btn btn-raised btn-danger btn-sm down-5px swal-demo2'><span class='ion-close-round'></span></button>"
					+ "</div>"
					+ "</div></div></div>";
				var li = $("<div class='col-md-3'></div>").html(html);
				li.appendTo('#listTag');
			}
		}
	});
}

function add(param) {
	var tag = {
		tag: param
	};

	$.ajax({
		url: '/v1/ui/tag/create?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(tag),
		success: function (data) {
			swal('Success!', 'Create success.');

			setTimeout(function () {
				window.location.reload();
			}, 1000);
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function deleteTag(param) {
	var tag = {
		tag: param
	};

	$.ajax({
		url: '/v1/ui/tag/delete?id=' + param + '&token=' + token,
		type: 'get',
		success: function (data) {
			location.reload();
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function deleteTagBulk(param) {

	var tag = {
		tags: param
	};

	$.ajax({
		url: '/v1/ui/tag/delete?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(tag),
		success: function (data) {
			location.reload();
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

(function () {
	'use strict';

	$(runSweetAlert);

	function runSweetAlert() {
		$(document).on('click', '.swal-demo4', function (e) {
			e.preventDefault();
			console.log(e.target.value);
			swal({
					title: 'Are you sure?',
					text: 'Do you want insert a new tag "' + $("#tag-value").val() + '"?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Insert',
					closeOnConfirm: false
				},
				function () {
					add($("#tag-value").val());
				});
		});
		$(document).on('click', '.swal-demo2', function (e) {
			e.preventDefault();
			console.log(e.target.value);
			swal({
					title: 'Are you sure?',
					text: 'Do you want delete tag "' + e.target.value + '"?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Delete',
					closeOnConfirm: false
				},
				function () {
					deleteTag(e.target.value);
				});
		});
		$(document).on('click', '.swal-demo3', function (e) {
			var checkbox = $("input[type=checkbox]:checked");
			var data = [];

			for (var i = 0; i < checkbox.length; i++) {
				data[i] = checkbox[i].value;
			}

			if(data.length == 0){
				swal('No tag selected.');
				return
			}

			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want delete all these tags?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Delete',
					closeOnConfirm: false
				},
				function () {
					deleteTagBulk(data);
				});
		});
	}

})();
