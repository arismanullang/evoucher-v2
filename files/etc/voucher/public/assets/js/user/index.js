function login(){
  var request = {
      username: $("#username").val(),
      password: $("#password").val()
  };

  $.ajax({
      url: 'http://evoucher.elys.id:8889/login',
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(request),
      success: function (data){
        alert(data.data);
        window.location = "http://evoucher.elys.id:8889/variant/create?user_id="+data.data;
      }
  });
}
