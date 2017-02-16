$( window ).load(function() {
  var token = findGetParameter("token");
  var id = findGetParameter("id");
  getName(id, token);
  searchByUser();
});

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

function toTwoDigit(val){
  if (val < 10){
    return '0'+val;
  }
  else {
    return val;
  }
}

function search() {
    var request = {
        fields: ["created_by", "variant_name"],
        values: ["IzKyd9yX", $("#variantName").val()]
      };

    $.ajax({
        url: '/variant/getVariant',
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify(request),
        success: function (data){
          renderData(data);
        }
    });
}

function searchByUser() {
    var request = {
        user: "IzKyd9yX" //$("#variantName").val()
      };

    $.ajax({
        url: 'http://evoucher.elys.id:8080/variant/getVariantByUser',
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify(request),
        success: function (data){
          renderData(data);
        }
    });
}

function getName(id, token) {
  $.get( 'http://juno-staging.elys.id:8888/v1/api/accounts/'+id+'?token='+token, function (data){
    $(".username").html(data.data.name);
  });
}

function login(username, password){
  $.ajax({
      url: 'http://juno-staging.elys.id:8888/v1/api/token',
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      beforeSend: function (xhr) {
        xhr.setRequestHeader ("Authorization", "Basic " + btoa(username+":"+password));
      },
      success: function (data){
        $.each(data, function(key, val) {
          $.each(val, function(k, v){
            if(k == "token"){
              alert(v);
            }
          });
        });

      }
  });
}

function renderData(data) {
  var arrData = [];

  $('#listProgram').html("");
  $.each(data, function(key, val) {
    $.each(val, function(k, v){
      if (k == "Data"){
        var length= v.length;

        $.each(v, function(x, y){
          var i = 0;
          var str = "";
          $.each(y, function(field, data){
            if(field != "CompanyID" && field != "ValidUsers"){
              arrData[i] = data;
              i++;
            }
          });

          var tr=$("<li class='li-first'></li>");
          var li = "";
          li = li + "<div class='row'>";
          li = li + "<div class='col-md-6'>Nama Program : <i>"+arrData[1]+"</i></div>";
          li = li + "<div class='col-md-6'><button type='button' class='btn btn-warning btn-xs' value="+arrData[0]+" onclick='goUpdate(this.value)'>Edit</button></div>";
          li = li + "</div>";
          $(tr).html(li);
          tr.appendTo('#listProgram');
          tr = $("<li></li>");

          var ul=$("<ul class='list-unstyled'></ul>");
          var li = "";
          li = li + "<li class='li-second'><div class='row'>";
          li = li + "<label class='col-md-6'>Voucher Type</label>";
          li = li + "<label class='col-md-6'>"+arrData[2]+"</label>";
          li = li + "</div></li>";
          li = li + "<li class='li-second'><div class='row'>";
          li = li + "<label class='col-md-6'>Harga Voucher</label>";
          li = li + "<label class='col-md-6'>"+arrData[3]+" Poin</label>";
          li = li + "</div></li>";
          li = li + "<li class='li-second'><div class='row'>";
          li = li + "<label class='col-md-6'>Limit Voucher</label>";
          li = li + "<label class='col-md-6'>"+arrData[4]+"</label>";
          li = li + "</div></li>";

          var d = new Date(arrData[5]);
          var date = toTwoDigit(d.getDate());
          var month = toTwoDigit(d.getMonth()+1);
          var hour = toTwoDigit(d.getHours());
          var minute = toTwoDigit(d.getMinutes());
          var second = toTwoDigit(d.getSeconds());
          var date = d.getFullYear()+'-'+month+'-'+date+' '+hour+':'+minute+':'+second;
          li = li + "<li class='li-second'><div class='row'>";
          li = li + "<label class='col-md-6'>Program Mulai</label>";
          li = li + "<label class='col-md-6'>"+date+"</label>";
          li = li + "</div></li>";

          d = new Date(arrData[6]);
          date = toTwoDigit(d.getDate());
          month = toTwoDigit(d.getMonth()+1);
          hour = toTwoDigit(d.getHours());
          minute = toTwoDigit(d.getMinutes());
          second = toTwoDigit(d.getSeconds());
          date = d.getFullYear()+'-'+month+'-'+date+' '+hour+':'+minute+':'+second;
          li = li + "<li class='li-second'><div class='row'>";
          li = li + "<label class='col-md-6'>Program Berakhir</label>";
          li = li + "<label class='col-md-6'>"+date+"</label>";
          li = li + "</div></li>";

          $(ul).html(li);
          $(tr).html(ul);
          tr.appendTo('#listProgram');
        });
      }
    });
  });
}

function goUpdate(data){
  window.location = "http://127.0.0.1:8080/variant/update?id="+data;
}
