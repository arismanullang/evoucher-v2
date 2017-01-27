$( window ).load(function() {
  searchByUser();
});

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
        values: ["nZ9Xmo-2", $("#variantName").val()]
      };

    $.ajax({
        url: 'http://127.0.0.1:8080/variant/getVariant',
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
        fields: "user",
        values: "a" //$("#variantName").val()
      };

    $.ajax({
        url: 'http://127.0.0.1:8080/variant/getVariantByUser',
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify(request),
        success: function (data){
          renderData(data);
        }
    });
}

function searchByName() {

    var arrData = [];
    var request = {
        fields: "variant_name",
        values: $("#variantName").val()
      };

    $.ajax({
      url: 'http://127.0.0.1:8080/variant/searchVariant',
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(request),
      success: function (data) {
        renderData(data);
      }
  });
}

function renderData(data) {
  var arrData = [];

  $('#listProgram').html("");
  $.each(data, function(key, val) {
    $.each(val, function(k, v){
      if (k == "VariantValue"){
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
          li = li + "<div class='col-md-6'><button type='button' class='btn btn-warning btn-xs'>Edit</button></div>";
          li = li + "</div>";
          $(tr).html(li);
          tr.appendTo('#listProgram');
          tr = $("<li></li>");

          var ul=$("<ul class='list-unstyled'></ul>");
          var li = "";
          li = li + "<li class='li-second'><div class='row'>";
          li = li + "<label class='col-md-6'>Point</label>";
          li = li + "<label class='col-md-6'>"+arrData[2]+"</label>";
          li = li + "</div></li>";
          li = li + "<li class='li-second'><div class='row'>";
          li = li + "<label class='col-md-6'>Limit Voucher</label>";
          li = li + "<label class='col-md-6'>"+arrData[3]+"</label>";
          li = li + "</div></li>";

          var d = new Date(arrData[4]);
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

          d = new Date(arrData[5]);
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
