var emailArr = [];
$(window).ready(function () {
	$("input").keydown(function(e) {
		if ( e.which == 13 ) e.preventDefault();
	});

	$("#search-email").on('keyup', function (e) {
		if (e.keyCode == 13) {
			searchEmail(this.value);
		}
	});

	$(".badge-remove").on('click', function (e) {
		var _this = $(this);
		_this.closest('.chosen').remove();
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

	var id = findGetParameter('id');
	$('#list-id').val(id);
	getUser(id);
});

function getUser(id) {
	$.ajax({
		url: '/v1/ui/user/list?token=' + token + '&id=' + id,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			$("#list-name").html(arrData.name);
			emailArr = arrData.email_users;
		}
	});
}

function onChangeElem(){
	$("#all-chosen-email").change(function () {
		var _this = $(this);
		_this.closest('#chosen-list').find("input.chosen").prop('checked', _this.prop('checked'));
	});
}

function send() {
	if(!$('#createListEmail').valid()){
		$(".error").focus();
		return;
	}

	var listEmail = [];
	var li = $(".chosen");

	$.each( li, function (i, val) {
		listEmail.push(val.getAttribute("value"));
	});

	var userReq = {
		email_user_id: listEmail,
		list_id: $("#list-id").val()
	};

	$.ajax({
		url: '/v1/ui/user/list/add-exist?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(userReq),
		success: function () {
			swal({
					title: 'Success',
					text: 'Submit Success',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/user/list/check?id="+$("#list-id").val();
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			if (a.errors.detail == "Duplicate Entry.") {
				swal("Username already used.");
			}
		}
	});
}

function searchEmail(param) {
	$.ajax({
		url: '/v1/ui/user/email?token='+token+'&email='+param,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			for (i = 0; i < arrData.length; i++) {
				var tc = checkCheckBox(arrData[i].id);
				if(tc){
					$('#box-list').attr('style', 'display:block');

					var li = $("<div class='col-sm-12'></div>");
					var html = "<label class='checkbox-inline c-checkbox'>"
						+ "<input type='checkbox' class='search' name='search' value='" + arrData[i].id + "' text='" + arrData[i].email;
					var ta = checkArr(arrData[i].id);
					if(ta){
						html += "' disabled><span class='ion-checkmark-round'></span>" + arrData[i].email
								+"<i class='text-danger'>   This email already in list</i>";
					}else{
						html += "''><span class='ion-checkmark-round'></span>" + arrData[i].email;
					}

					html += "</label>";
					li.html(html);
					li.appendTo('#search-list');
				}
			}

			$(".search").change(function () {
				var _this = $(this);
				addCheckBox(_this.val(), _this.attr('text'), 'chosen', 'chosen', 'chosen-list');
				_this.parent().closest('.col-sm-12').remove();

				var tf = checkSearch();
				if(!tf){
					$('#box-list').attr('style', 'display:none');
				}
			});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function checkSearch(){
	var tf = false;
	var li = $("input[class=search]");

	if(li.length > 0){
		tf = true;
	}
	return tf;
}

function checkCheckBox(param){
	var li = $(".chosen");
	var tf = true;

	$.each( li, function (i, val) {
		if(val.getAttribute("value") == param){
			tf = false;
		}
	});

	return tf;
}

function checkArr(param){
	tf = false;
	if(emailArr != null){
		for(var i = 0; i < emailArr.length; i++){
			if(emailArr[i].id == param){
				tf = true;
			}
		}
	}

	return tf;
}

function removeCheckBox(param){
	var li = $("input[class="+param+"]:checked");

	$.each( li, function (i, val) {
		$(val).parent().closest('.col-sm-12').remove();
	});
}

function addCheckBox(id, email, classElem, nameElem, listElem){
	var li = $("<span class='d-inline p-2 bg-divo m-sm p box "+classElem+"' name='"+nameElem+"' value='"+id+"'></span>");
	var html = "<span class='badge bg-white-grey mr-sm badge-remove'>X</span>"
		+ email;

	li.html(html);
	li.appendTo("#"+listElem);

	$(".badge-remove").on('click', function (e) {
		var _this = $(this);
		_this.closest('.chosen').remove();
	});
}
