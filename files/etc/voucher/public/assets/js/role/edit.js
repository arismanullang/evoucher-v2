$( document ).ready(function() {
	id = findGetParameter("id");
	getFeature(id);
	$("#role-id").val(id);

	$("#all-ui-feature").change(function() {
		var _this = $(this);
		_this.closest('.checked-container').find("input[class=feature]").prop('checked',_this.prop('checked'));
	});

	$("#all-api-feature").change(function() {
		var _this = $(this);
		_this.closest('.checked-container').find("input[class=feature]").prop('checked',_this.prop('checked'));
	});
});

function getFeature(id) {
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
				if(arrData[i].type == 'ui'){
					uiFeature.push(arrData[i]);
				}else{
					apiFeature.push(arrData[i]);
				}
			}

			generateElem(uiFeature, '#list-ui-features');
			generateElem(apiFeature, '#list-api-features');

			getFeatureDetail(id);
		}
	});
}

function getFeatureDetail(id) {
	$.ajax({
		url: '/v1/ui/role/detail?id='+id+'&token='+token,
		type: 'get',
		success: function (data) {
			var arrData = data.data;
			$("#role-detail").html(arrData.detail);
			if(arrData.features != null){
				for (i = 0; i < arrData.features.length; i++){
					$("#features").find("input[id="+arrData.features[i]+"]").prop('checked', true);
				}
			}
		},
		error: function (data) {
			var a = JSON.parse(data.responseText);
			swal("Error", a.errors.detail);
		}
	});
}

function update() {
	var uiElems = $("#list-ui-features").find("input[class=feature]:checked");
	var uis = [];

	for(var i = 0; i < uiElems.length; i++){
		uis.push($(uiElems[i]).attr("feature"));
	}
	localStorage.removeItem("ui");
	localStorage.setItem("ui", uis);

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

function generateElem(param, targetElem){
	var categories = [];
	for(i = 0; i < param.length; i++){
		if(i == 0){
			categories.push(param[i].category);
		}else{
			var tf = true;
			for(y = 0; y < categories.length;y++){
				if(param[i].category == categories[y]){
					tf = false;
				}
			}

			if(tf){
				categories.push(param[i].category);
			}
		}
	}

	for(i = 0; i < categories.length/2; i++){
		var ii = i*2;
		var html = "<div class='col-sm-6 display-flex'><div class='card card100'><div class='card-heading heading-sub'>"
			+ "<div class='card-title'>"+categories[ii]+"</div>"
			+ "</div>"
			+ "<div class='card-body pt0 pb0'>"
			+ "<div class='row'>";

		for(y = 0; y < param.length; y++){
			if(categories[ii] == param[y].category){
				html +="<div class='col-sm-6'>"
				+ "<div class='checkbox c-checkbox'>"
				+ "<label class='text-thin font-size-12px'>"
				+ "<input id='" + param[y].id + "' value='" + param[y].id + "' type='checkbox' class='feature' feature='/"+categories[ii]+"/"+param[y].detail+"'><span class='ion-checkmark-round'></span>" + param[y].detail
				+ "</label>"
				+ "</div>"
				+ "</div>";
			}
		}

		html += "</div></div></div></div>";

		if((ii+1) < categories.length){
			html += "<div class='col-sm-6 display-flex'><div class='card card100'><div class='card-heading heading-sub'>"
				+ "<div class='card-title'>"+categories[ii+1]+"</div>"
				+ "</div>"
				+ "<div class='card-body pt0 pb0'>"
				+ "<div class='row'>";

			for(y = 0; y < param.length; y++){
				if(categories[ii+1] == param[y].category){
					html +="<div class='col-sm-6'>"
					+ "<div class='checkbox c-checkbox'>"
					+ "<label class='text-thin font-size-12px'>"
					+ "<input id='" + param[y].id + "' value='" + param[y].id + "' type='checkbox' class='feature' feature='/"+categories[ii+1]+"/"+param[y].detail+"'><span class='ion-checkmark-round'></span>" + param[y].detail
					+ "</label>"
					+ "</div>"
					+ "</div>";
				}
			}

			html += "</div></div></div></div>";
		}
		var li = $("<div class='row table-row'></div>").html(html);
		li.appendTo(targetElem);
	}
}
