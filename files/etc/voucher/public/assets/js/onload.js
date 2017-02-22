$( window ).load(function() {
  var token = findGetParameter("token");
  getSession(token);
});

function getSession(token) {
    $.ajax({
      url: 'http://evoucher.elys.id:8889/get/session?token='+token,
      type: 'get',
      success: function (data) {
        alert(data.data);
        if(data.data == false){
          window.location = "http://evoucher.elys.id:8889/user/login";
        }
      }
  });
}

function findGetParameter(parameterName) {
    var result = null,
        tmp = [];
    location.search
    .substr(1)
        .split("&")
        .forEach(function (item) {
        tmp = item.split("=");
        if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
    });
    return result;
}
