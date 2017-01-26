function toTwoDigit(val){
  if (val < 10){
    return '0'+val;
  }
  else {
    return val;
  }
}

function send() {
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
        validUsers: ["nZ9Xmo-2", "nZ9Xmo-2"]
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

function search() {

    var arrData = [];
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
        success: function (data) {
          $('#searchTable').html("");
          var th = "<tr>"+
            "<th>Point</th>"+
            "<th>Name</th>"+
            "<th>Max Vouchers</th>"+
            "<th>Start</th>"+
            "<th>End</th>"+
            "<th>T&C</th>"+
            "</tr>";

          $(th).appendTo('#searchTable');

          $.each(data, function(key, val) {
            $.each(val, function(k, v){
                if (k == "VariantValue"){
                  $.each(v, function(x, y){
                    var tr=$('<tr></tr>');
                    var i = 0;
                    $.each(y, function(field, data){
                      if(field != "CompanyID" && field != "ValidUsers"){
                        arrData[i] = data;
                        i++;
                      }
                    });

                    for (var z = 1; z < arrData.length; z++) {
                      if( z == 4 || z == 5){
                        var d = new Date(arrData[z]);
                        var date = toTwoDigit(d.getDate());
                        var month = toTwoDigit(d.getMonth()+1);
                        var hour = toTwoDigit(d.getHours());
                        var minute = toTwoDigit(d.getMinutes());
                        var second = toTwoDigit(d.getSeconds());

                        var date = d.getFullYear()+'-'+month+'-'+date+' '+hour+':'+minute+':'+second;
                        $('<td>'+date+'</td>').appendTo(tr);
                      }
                      else{
                        $('<td>'+arrData[z]+'</td>').appendTo(tr);
                      }
                    }
                    tr.appendTo('#searchTable');
                  });
                }
            });
          });
        }
    });
}
