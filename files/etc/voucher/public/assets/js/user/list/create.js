$(window).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	$("#search-email").on('keyup', function (e) {
		if (e.keyCode == 13) {
			searchEmail(this.value);
		}
	});

	$('#createListEmail').validate({
		errorPlacement: errorPlacementInput,
		// Form rules
		rules: {
			name: {
				required: true
			}
		}
	});

	onChangeElem();
});

function onChangeElem(){
	$("#all-email").change(function () {
		var _this = $(this);
		_this.closest('#email-list').find("input.email").prop('checked', _this.prop('checked'));
	});
}

function send() {
	if($("#name").val() == ''){
	    var errors = { name: "List Email Cannot be Empty." };
        /* Show errors on the form */
        $("#createListEmail").validate().showErrors(errors);
		$(".error").focus();
		return;
	}

	var listEmail = [];
	var li = $("input[class=email]:checked");

	$.each( li, function (i, val) {
		listEmail.push(val.getAttribute("value"));
	});

	var userReq = {
		name: $("#name").val(),
		email_users: listEmail
	};

	console.log(userReq);
	$.ajax({
		url: '/v1/ui/user/list/create?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(userReq),
		success: function () {
			swal({
					title: 'Success',
					text: 'Create List Success',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/user/list/search";
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal(a.errors.detail);
		}
	});
}

function searchEmail(param) {
	$.ajax({
		url: '/v1/ui/user/email?token='+token+'&email='+param,
		type: 'get',
		success: function (data) {
			var liChecked = $("input[class=email]:checked");
			var liNotChecked = $("input[class=email]:not(:checked)");

			$.each( liNotChecked, function (i, val) {
				$(val).parent().closest('.col-sm-4').remove();
			});

			var arrData = [];
			arrData = data.data;
			var i;
			for (i = 0; i < arrData.length; i++) {
				var tf = true;
				for(var y = 0; y < liChecked.length;y++){
					if($(liChecked[y]).val() == arrData[i].id){
						tf = false;
					}
				}

				if(tf){
					var li = $("<div class='col-sm-4'></div>");
					var html = "<label class='checkbox-inline c-checkbox'>"
						+ "<input type='checkbox' name='email' class='email' value='" + arrData[i].id + "'>"
						+ "<span class='ion-checkmark-round'></span>" + arrData[i].email
						+ "</label>";
					li.html(html);
					li.appendTo('#email-list');
				}
			}
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
