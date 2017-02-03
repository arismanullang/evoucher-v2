$( window ).load(function() {
  findGetParameter("id");
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
    searchById(result);
}

function searchById(id) {

    var arrData = [];

    $.ajax({
        url: 'http://127.0.0.1:8080/variant/'+id,
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        success: function (data) {
          $.each(data, function(key, val) {
            $.each(val, function(k, v){
                if (k == "Data"){
                  $.each(v, function(x, y){
                    switch(x){
                      case "ID":
                        $("#variantId").val(y);
                        break;
                      case "VariantName":
                        $("#variantName").val(y);
                        break;
                      case "VoucherType":
                        $("#voucherType").val(y);
                        break;
                      case "VariantType":
                        $("#variantType").val(y);
                        break;
                      case "DiscountValue":
                        $("#discountValue").val(y);
                        break;
                      case "MaxVoucher":
                        $("#maxVoucher").val(y);
                        break;
                      case "StartDate":
                        $("#startDate").val(y);
                        break;
                      case "EndDate":
                        $("#endDate").val(y);
                        break;
                      case "AllowAccumulative":
                        $("#allowAccumulative").prop('checked', y);
                        break;
                      case "PointNeeded":
                        $("#pointNeeded").val(y);
                        break;
                      case "ImgUrl":
                        $("#imgUrl").val(y);
                        break;
                      case "VariantTnc":
                        $("#variantTnc").val(y);
                        break;

                    }
                  });
                }
            });
          });
        }
    });
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
        url: 'http://127.0.0.1:8080/variant/'+$("#variantId").val()+'/update',
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify(variant),
        success: function () {
            alert("Variant Edited.");
        }
    });
}
