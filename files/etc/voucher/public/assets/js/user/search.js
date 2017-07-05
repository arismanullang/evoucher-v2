$( document ).ready(function() {
	getUser();
});

function getUser() {
	$.ajax({
		url: '/v1/ui/user/all?token='+token,
		type: 'get',
		success: function (data) {
			console.log("Render Data");
			var arrData = [];
			arrData = data.data;
			console.log(arrData);
			var i;
			var dataSet = [];
			for (i = 0; i < arrData.length; i++){
				var button = "<button type='button' class='btn btn-flat btn-sm btn-info' onclick='edit(\""+arrData[i].id+"\")'><em class='ion-edit'></em></button>"+
					"<button value='"+arrData[i].id+"' type='button' class='btn btn-flat btn-sm btn-danger swal-demo-reset'><em class='ion-locked'></em></button>"+
					"<button value='"+arrData[i].id+"' type='button' status='"+arrData[i].status+"' class='btn btn-flat btn-sm btn-danger swal-demo-delete'><em class='ion-trash-a'></em></button>";
				var status = "INACTIVE";
				if(arrData[i].status == "created"){
					status = "ACTIVE";
				}
				dataSet[i] = [
					arrData[i].username
					, arrData[i].email
					, arrData[i].phone
					, status
					, button
				];
			}

			if ($.fn.DataTable.isDataTable("#datatable1")) {
				$('#datatable1').DataTable().clear().destroy();
			}

			var table = $('#datatable1').dataTable({
				data: dataSet,
				dom: 'rtip',
				"order": [[ 1, "desc" ]],
				columns: [
					{ title: "USERNAME" },
					{ title: "EMAIL" },
					{ title: "PHONE" },
					{ title: "STATUS"},
					{ title: "ACTION"}
				],
				oLanguage: {
					sSearch: '<em class="ion-search"></em>',
					sLengthMenu: '_MENU_ records per page',
					info: 'Showing page _PAGE_ of _PAGES_',
					zeroRecords: 'Nothing found - sorry',
					infoEmpty: 'No records available',
					infoFiltered: '(filtered from _MAX_ total records)',
					oPaginate: {
						sNext: '<em class="ion-ios-arrow-right"></em>',
						sPrevious: '<em class="ion-ios-arrow-left"></em>'
					}
				}
			});
			var inputSearchClass = 'datatable_input_col_search';
			var columnInputs = $('thead .' + inputSearchClass);

			columnInputs
				.keyup(function() {
					table.fnFilter(this.value, columnInputs.index(this));
				});
		}
	});
}

function edit(url){
	window.location = "/user/update?id="+url+"&token="+token;
}

function addUser() {
	window.location = "/user/register?token="+token;
}

function resetPassword(id) {
	var newPass = "";
	var user = {
		id: id
	};
	$.ajax({
		url: '/v1/ui/user/update?type=reset&token='+token,
		type: 'POST',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			newPass = data.data;
			console.log('Reset : '+newPass);
			swal('Success!', "Reset success. \n New Password : "+newPass, getUser());
		}
	});
}

function deleteUser(id) {

	var user = {
		id: id
	};
	$.ajax({
		url: '/v1/ui/user/block?token='+token,
		type: 'POST',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			getUser();
		}
	});
}

function activateUser(id) {
	var user = {
		id: id
	};
	$.ajax({
		url: '/v1/ui/user/activate?token='+token,
		type: 'POST',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(user),
		success: function (data) {
			getUser();
		}
	});
}

(function() {
	'use strict';

	$(runSweetAlert);
	function runSweetAlert() {
		$(document).on('click', '.swal-demo-delete', function(e) {
			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want change user\'s status?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Change!',
					closeOnConfirm: false
				},
				function() {
					if(e.target.getAttribute("status") == "created"){
						swal('Blocked!', 'Blocked success.', deleteUser(e.target.value));
					}else {
						swal('Activated!', 'Activated success.', activateUser(e.target.value));
					}

				});

		});
		$(document).on('click', '.swal-demo-reset', function(e) {
			e.preventDefault();
			swal({
					title: 'Are you sure?',
					text: 'Do you want reset user\'s password?',
					type: 'warning',
					showCancelButton: true,
					confirmButtonColor: '#DD6B55',
					confirmButtonText: 'Reset!',
					closeOnConfirm: false
				},
				function() {
					resetPassword(e.target.value);
				});

		});
	}

})();
