//Usage
$( window ).ready(function() {
	  	$(".toast").attr("style","display:none");
		$(".modal-qr").modal();
		localStorage.removeItem('emailSpin');
		localStorage.removeItem('nameSpin');
});

var token = localStorage.getItem("token");

//load your JSON (you could jQuery if you prefer)
function loadJSON(callback) {
  var xobj = new XMLHttpRequest();
  xobj.overrideMimeType("application/json");
  xobj.open('GET', '/v1/ui/program/spin?token='+token, true);
  xobj.onreadystatechange = function() {
    if (xobj.readyState == 4 && xobj.status == "200") {
      //Call the anonymous function (callback) passing in the response
	  callback(xobj.responseText);
    }
  };
  xobj.send(null);
}

//your own function to capture the spin results
function myResult(e) {
	$(".arrow").attr("style","display:none");
  	$(".toast").css("display","block");
  	//e is the result object
 	//console.log(e);
    // console.log('Spin Count: ' + e.spinCount + ' - ' + 'Win: ' + e.win + ' - ' + 'Message: ' +  e.msg.split("`")[1]);

	var response = e.msg.split("`");
	var program = {
		program_id: response[2],
		reference_no: e.gameId,
		holder: {
			id: localStorage.getItem("emailSpin"),
			email: localStorage.getItem("emailSpin"),
			phone: "",
			description: localStorage.getItem("nameSpin")
		},
		subject: "Gift Voucher "+response[1]
	};

	$.ajax({
		url: '/v1/ui/voucher/generate/email?token=' + token,
		type: 'post',
		dataType: 'json',
		contentType: "application/json",
		data: JSON.stringify(program),
		success: function () {

		},
		error: function (data) {
			myError(e);
		}
	});
}

//your own function to capture any errors
function myError(e) {
  //e is error object
  console.log('Spin Count: ' + e.spinCount + ' - ' + 'Message: ' +  e.msg);

}

function myGameEnd(e) {

  //e is gameResultsArray
  console.log(e);
  TweenMax.delayedCall(6, function(){
//    window.location.reload();
  })


}

function init() {
  loadJSON(function(response) {
    // Parse JSON string to an object
    var jsonData = JSON.parse(response);
    //if you want to spin it using your own button, then create a reference and pass it in as spinTrigger
    var mySpinBtn = document.querySelector('.spinBtn');
    //create a new instance of Spin2Win Wheel and pass in the vars object
    var myWheel = new Spin2WinWheel();

    //WITH your own button
    myWheel.init({data:jsonData.data, onResult:myResult, onGameEnd:myGameEnd, onError:myError, spinTrigger:mySpinBtn});

    //WITHOUT your own button
    //myWheel.init({data:jsonData, onResult:myResult, onGameEnd:myGameEnd, onError:myError});
  });
}

function add(){
	localStorage.setItem("userSpin", $("#name").val());
	localStorage.setItem("emailSpin", $("#email").val());
}

//And finally call it
init();
