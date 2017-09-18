$( document ).ready(function() {
	id = findGetParameter("id");
	getFeature();
	getFeatureDetail(id);

	$("#all-ui-feature").change(function() {
		var _this = $(this);
		_this.closest('.checked-container').find("input[class=feature]").prop('checked',_this.prop('checked'));
	});

	$("#all-api-feature").change(function() {
		var _this = $(this);
		_this.closest('.checked-container').find("input[class=feature]").prop('checked',_this.prop('checked'));
	});
});

function getFeature() {
	$.ajax({
		url: '/v1/ui/feature/all?token='+token,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			var uiFeature = [];
			var apiFeature = [];
			for (i = 0; i < arrData.length; i++){
				var html = "<div class='card'><div class='card-body pt0 pb0'><div class='row'>"
					+ "<div class='col-sm-12'>"
					+ "<div class='checkbox c-checkbox'>"
					+ "<label class='text-thin font-size-12px'>"
					+ "<input id='"+arrData[i].id+"' value='"+arrData[i].id+"' type='checkbox' class='feature'><span class='ion-checkmark-round'></span>"+arrData[i].category+"-"+arrData[i].detail
					+ "</label>"
					+ "</div>"
					+ "</div>"
					+ "</div></div></div>";
				var li = $("<div class='col-md-6'></div>").html(html);

				if(arrData[i].type == 'ui'){
					li.appendTo('#list-ui-feature');
				}else{
					li.appendTo('#list-api-feature');
				}
			}
		}
	});
}

function getFeatureDetail(id) {
	$.ajax({
		url: '/v1/ui/role/detail?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var arrData = data.data;
			$("#role-id").val(arrData.id);
			$("#role-detail").html(arrData.detail);
			for (i = 0; i < arrData.features.length; i++){
				$("#features").find("input[id="+arrData.features[i]+"]").prop('checked', true);
			}
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function update() {

	var listFeatures = [];
	var li = $( "input[class=feature]:checked" );

	if(li.length == 0 || parseInt($("#length").val()) < 8){
		error = true;
		errorMessage = "Select Feature. ";
	}

	for (i = 0; i < li.length; i++) {
		if(li[i].value != "on") {
			listFeatures[i] = li[i].value;
		}
	}

	var role = {
		id: $("#role-id").val(),
		features: listFeatures
	};

	$.ajax({
		url: '/v1/ui/role/update?token='+token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(role),
		success: function (data) {
			swal({
					title: 'Success',
					text: 'Role Updated',
					type: 'success',
					showCancelButton: false,
					confirmButtonText: 'Ok',
					closeOnConfirm: false
				},
				function() {
					window.location = "/role/search";
				});
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}
