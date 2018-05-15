$(document).ready(function () {
	var id = findGetParameter("id");
	$("#program-id").val(id);
	getUser();
	localStorage.removeItem('list-email-id');
});

function getUser() {
	$.ajax({
		url: '/v1/ui/user/list/all?token=' + token,
		type: 'get',
		success: function (data) {
			var arrData = [];
			arrData = data.data;
			var i;
			var dataSet = [];

			for(var i = 0; i < arrData.length; i++){
				var emailUsers = 0;
				if(arrData[i].email_users != null){
					emailUsers = arrData[i].email_users.length;
				}
				var body = "<td class='col-lg-1 checkbox c-checkbox'><label>"
					+ "<input type='checkbox' name='email' class='email' value='"+arrData[i].id+"'><span class='ion-checkmark-round'></span>"
					+ "</label></td>"
					+ "<td class='text-ellipsis'>"+arrData[i].name+"</td>"
					+ "<td class='text-ellipsis'>"+emailUsers+"</td>"
				var li = $("<tr></tr>");
				li.html(body);
				li.appendTo('#list-email');
			}
		}
	});
}

function next() {
	var li = $("input[class=email]:checked");
	var listEmailId = li[0].value;
	if( li.length > 1){
		for(var i = 1; i < li.length; i++){
			listEmailId += "`"+li[i].value;
		}
	}
	localStorage.setItem("list-email-id", listEmailId);
 	window.location = "/program/send-campaign?id="+$("#program-id").val();
}
