$( window ).load(function() {
  searchByRole();
  sess();
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
  var listName = [];
  for (i = 0; i < $(".name-op").length; i++){
    listName[i] = $(".name-op").eq(i).text();
  }

  var variant = {
      companyId: $("#companyId").val(),
      variantName: $("#variantName").val(),
      variantType: $("#variantType").val(),
      pointNeeded: parseInt($("#pointNeeded").val()),
      maxVoucher: parseInt($("#maxVoucher").val()),
      allowAccumulative: $("#allowAccumulative").is(":checked"),
      startDate: $("#startDate").val(),
      finishDate: $("#endDate").val(),
      discountValue: parseInt($("#discountValue").val()),
      imgUrl: $("#imgUrl").val(),
      variantTnc: $("variantTnc").val(),
      createdBy: "nZ9Xmo-2",
      validUsers: listName
    };

  $.ajax({
      url: 'http://127.0.0.1:8080/variant/createVariant',
      type: 'post',
      dataType: 'json',
      contentType: "application/json",
      data: JSON.stringify(variant),
      success: function () {
          alert("Variant created.");
      }
  });
}

function sess() {
    $.ajax({
      url: 'http://evoucher.elys.id:8889/get/session',
      type: 'get',
      success: function (data) {
        alert(data);
      }
  });
}

function searchByRole() {

    var arrData = [];
    var request = {
        role: "operator"
      };

    $.ajax({
      url: 'http://127.0.0.1:8080/user/getUserByRole',
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

  $('#listOperator').html("");
  $.each(data, function(key, val) {
    $.each(val, function(k, v){
      if (k == "Data"){
        var length= v.length;

        $.each(v, function(x, y){
          var i = 0;
          var str = "";
          $.each(y, function(field, data){
            arrData[i] = data;
            i++;
          });
          var tr=$("<li></li>");
          var li = "";
          li = li + "<button type='button' class='btn btn-list' value="+arrData[0]+" onclick='addData(this)'>"+arrData[1]+"</button>";
          $(tr).html(li);
          tr.appendTo('#listOperator');
        });
      }
    });
  });
}

function addData(elem) {
  var tr=$("<span class='label label-success name-op' onclick='remove(this)'></span>");
  $(tr).html(elem.value);
  tr.appendTo('#listOp');

}

function remove(elem){
  $(elem).remove();
}
