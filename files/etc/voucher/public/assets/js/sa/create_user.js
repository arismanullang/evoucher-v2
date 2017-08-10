$( window ).ready(function() {
  getRole();
  getAccount();
});

function toTwoDigit(val){
  if (val < 10){
    return '0'+val;
  }
  else {
    return val;
  }
}

function send() {
  var i;

  var listRole = [];
  var li = $( "input[type=checkbox]:checked" );

  if(li.length == 0 || parseInt($("#length").val()) < 8){
    error = true;
  }

  for (i = 0; i < li.length; i++) {
    listRole[i] = li[i].value;
  }

  var error = false;
  $('input[check="true"]').each(function() {
    if($(this).val() == ""){
      $(this).addClass("error");
      $(this).parent().closest('div').addClass("input-error");
      error = true;
    }
  });

  if(error){
    alert("Please check your input.");
    return
  }

  var userReq = {
      account_id: $("#account").val(),
      username: $("#username").val(),
      password: $("#password").val(),
      email: $("#email").val(),
      phone: $("#phone").val(),
      role_id: listRole
    };

    console.log(userReq);
    $.ajax({
       url: '/v1/ui/sa/create?token='+token,
       type: 'post',
       dataType: 'json',
       contentType: "application/json",
       data: JSON.stringify(userReq),
       success: function () {
           alert("User created.");
           // window.location = "/sa/search?token="+token;
       },
    	error: function (data) {
		var a = JSON.parse(data.responseText);
		if(a.errors.detail == "Duplicate Entry."){
			alert("Username already used.");
		}
    	}
   });
}

function getRole() {
    console.log("Get Role");

    $.ajax({
      url: '/v1/ui/role/all',
      type: 'get',
      success: function (data) {
        var arrData = [];
        arrData = data.data;

        var i;
        for (i = 0; i < arrData.length; i++){
		var li = $("<div class='col-sm-4'></div>");
		var html = "<label class='checkbox-inline c-checkbox'>"
			+ "<input type='checkbox' value='"+arrData[i].id+"' text='"+arrData[i].detail+"'>"
			+ "<span class='ion-checkmark-round'></span>" + arrData[i].detail
			+ "</label>";
		li.html(html);
		li.appendTo('#role');
        }
      }
  });
}

function getAccount() {
	console.log("Get Account");

	$.ajax({
		url: '/v1/ui/account/all',
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			console.log(data.data);
			var i;
			for (i = 0; i < arrData.length; i++){
				// <option value="bulk">Email Blast</option>
				var li = $("<option value='"+arrData[i].id+"'></option>");
				li.html(arrData[i].name);
				li.appendTo('#account');
			}

			$('.select2').select2();
		}
	});
}
