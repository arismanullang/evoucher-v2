var id = findGetParameter("id");
var user = localStorage.getItem("user");
var token = localStorage.getItem(user);

$( window ).ready(function() {
  searchById(id);
  // getPartner();
});

function searchById(id) {

    var arrData = [];

    $.ajax({
        url: '/v1/api/get/variant/'+id+'?token='+token,
        type: 'get',
        success: function (data) {
          console.log(data);
          var variant = data.data;

          $("#variant-name").val(variant.VariantName);
          $("#variant-type").val(variant.VariantType);
          $("#select2-variant-type-container").html(convcertToUpperCase(variant.VariantType));
          $("#voucher-type").val(variant.VoucherType);
          $("#select2-voucher-type-container").html(convcertToUpperCase(variant.VoucherType));
          $("#voucher-price").val(variant.VoucherPrice);
          $("#max-quantity-voucher").val(variant.MaxQuantityVoucher);
          $("#max-usage-voucher").val(variant.MaxUsageVoucher);
          $("#redemption-method").val(variant.RedeemtionMethod);
          $("#select2-redemption-method-container").html(convcertToUpperCase(variant.RedeemtionMethod));
          $("#start-date").val(convertToDate(variant.StartDate));
          $("#end-date").val(convertToDate(variant.EndDate));
          $("#voucher-value").val(variant.DiscountValue);
          $("#image-url").val(variant.ImgUrl);
          $("#variant-tnc").val(variant.VariantTnc);
          $("#variant-description").val(variant.VariantDescription);

        }
    });
}

function send() {
    $('input[check="true"]').each(function() {
      if($(this).val() == ""){
        $(this).addClass("error");
        $(this).parent().closest('div').addClass("input-error");
        error = true;
      }

      if($(this).attr("id") == "length"){
        if(parseInt($(this).val()) < 8){
          error = true;
        }
      }
    });

    if(error){
      alert("Please check your input.");
      return
    }

    var variant = {
      variant_name: $("#variant-name").val(),
      variant_type: $("#variant-type").val(),
      voucher_type: $("#voucher-type").val(),
      voucher_price: parseInt($("#voucher-price").val()),
      max_quantity_voucher: parseInt($("#max-quantity-voucher").val()),
      max_usage_voucher: parseInt($("#max-usage-voucher").val()),
      redeemtion_method: $("#redemption-method").val(),
      start_date: $("#start-date").val(),
      end_date: $("#end-date").val(),
      discount_value: parseInt($("#voucher-value").val()),
      image_url: $("#image-url").val(),
      variant_tnc: $("#variant-tnc").val(),
      variant_description: $("#variant-description").val()
    };

    console.log(variant);

    $.ajax({
        url: '/v1/update/variant/'+id+'?token='+token,
        type: 'post',
        dataType: 'json',
        contentType: "application/json",
        data: JSON.stringify(variant),
        success: function () {
            window.location = "/variant/search";
        }
    });
}

function getPartner() {
    console.log("Get Partner Data");

    $.ajax({
      url: '/v1/get/partner',
      type: 'get',
      success: function (data) {
        console.log("Render Data");
        var arrData = [];
        arrData = data.data;

        var i;
        for (i = 0; i < arrData.length; i++){
          var li = $("<option value='"+arrData[i].Id+"'>"+arrData[i].PartnerName+"</option>");
          li.appendTo('#variantPartners');
        }
      }
  });
}

function convertToDate(date){
  var string1 = date.split("T")[0];
  var string2 = string1.split("-");
  var result = string2[1] + "/" + string2[2] + "/" + string2[0];

  return result;
}

function convcertToUpperCase(upperCase){
  var result = "";
  var firstChar = upperCase.charAt(0);
  upperCase = upperCase.replace(firstChar, firstChar.toUpperCase());
  result = upperCase;

  return result;
}

(function() {
    'use strict';

    $(formAdvanced);

    function formAdvanced() {
        $('.select2').select2();

        $('.datepicker4')
            .datepicker({
                container:'#example-datepicker-container-4'
            });
    }

})();
