$( document ).ready(function() {
	var total = 0;
	$('#transaction').submit(function(e) {
		e.preventDefault();
		return false;
	});

	var id = findGetParameter('id');
	getCampaign(id);
	getEmail(localStorage.getItem('list-email-id'));
	getListEmail(localStorage.getItem('list-email-id'));
});
var email = [];
function getCampaign(id) {
	$.ajax({
		url: '/v2/ui/campaign?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;

			$('#campaign-name').html(result.program_name);
			$('#subject-name').html(result.email_subject);
			$('#email-sender').html(result.email_sender);
		},
		error: function (data) {
		}
	});
}

function getEmail(id) {
	$.ajax({
		url: '/v1/ui/user/email-id?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;

			$('#total-email').html(result.length + " people");
			for(var i = 0; i < result.length;i++){
				email.push(result[i].id);
			}
		},
		error: function (data) {
		}
	});
}

function getListEmail(id) {
	$.ajax({
		url: '/v1/ui/user/lists?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var result = data.data;
			var name = result[0].name;
			for(var i = 1; i < result.length; i++){
				name += ", " + result[i].name;
			}
			$('#list-name').html(name);
		},
		error: function (data) {
		}
	});
}

function send(){
	var id = findGetParameter('id');

	var param = {
		program_id : id,
		email_user : email
	};

	$.ajax({
		url: '/v2/ui/campaign/send?token='+token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(param),
		success: function (data) {
			console.log(data);
			swal("Success", data);
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
